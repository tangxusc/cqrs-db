package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type AggregateCache struct {
	UpdateTime *time.Time
	Data       string
	Version    int
}

func NewAggregateCache() *AggregateCache {
	return &AggregateCache{
		UpdateTime: nil,
		Data:       "",
		Version:    0,
	}
}

/*
聚合
*/
type Aggregate struct {
	AggId        string
	AggType      string
	eventsChan   chan Events
	recoveryChan chan Events
	ctx          context.Context

	cache *AggregateCache
}

func NewAggregate(aggId, aggType string, ctx context.Context) (*Aggregate, error) {
	agg := &Aggregate{
		AggId:        aggId,
		AggType:      aggType,
		eventsChan:   make(chan Events, 10),
		recoveryChan: make(chan Events, 10),
		cache:        NewAggregateCache(),
	}
	e := agg.SyncCache()
	if e != nil {
		return nil, e
	}
	agg.Start(ctx)

	return agg, nil
}

func (a *Aggregate) Start(ctx context.Context) {
	a.ctx = ctx
	for {
		select {
		//优先级1
		case <-ctx.Done():
			close(a.eventsChan)
			close(a.recoveryChan)
			return
		default:
			var events Events
			select {
			//优先2发送recovery
			case events = <-a.recoveryChan:
				a.sendEvent(events)
			default:
				select {
				//优先级3
				case events = <-a.eventsChan:
					a.sendEvent(events)
				default:
					break
				}
			}
		}
	}
}

func (a *Aggregate) sendEvent(events Events) {
	var event *Event
	var e error
	for i := 0; i < len(events); {
		event = events[i]
		e = eventSender.Send(event)
		if e != nil {
			continue
		}
		i = i + 1
		//发送成功,但是数据库写入出错了
		//TODO:这个错误怎么办? 是否可以放到recovery中?
		e = event.SuccessSend()
	}
}

func (a *Aggregate) PutSendChan(events Events) error {
	//保存event
	err := eventRepository.SaveEvents(events)
	if err != nil {
		return err
	}
	select {
	case <-a.ctx.Done():
		return fmt.Errorf("已关闭")
	default:
		select {
		case a.eventsChan <- events:
			//TODO:考虑并发
			//计算快照,并缓存data
			e := a.applyEvents(a.cache.Data, events)
			if e != nil {
				return e
			}
			go func() {
				allow := snapshotSaveStrategy.Allow(a.AggId, a.AggType, a.cache.Data, events)
				if allow {
					snapshotStore.Save(a.AggId, a.AggType, a.cache)
				}
			}()
			return nil
		default:
			return fmt.Errorf("队列已满")
		}
	}
}

func (a *Aggregate) PutRecoveryChan(events Events) error {
	select {
	case <-a.ctx.Done():
		return fmt.Errorf("已关闭")
	default:
		select {
		case a.recoveryChan <- events:
			return nil
		default:
			return fmt.Errorf("队列已满")
		}
	}
}

/*
1,查询快照
2,根据快照查询事件
3,进行json merge溯源
4,返回结果
*/
func (a *Aggregate) Sourcing() (string, error) {
	return a.cache.Data, nil
}

func Sourcing(id, aggType string) (data string, e error) {
	agg := aggregateCache.Get(id, aggType)
	if agg == nil {
		return "", nil
	}
	return agg.Sourcing()
}

func (a *Aggregate) applyEvents(data string, events []*Event) (e error) {
	if len(events) == 0 {
		return
	}
	//按照顺序合并数据
	//聚合版本字段处理,如果小于当前版本,则需要丢弃这个event
	for _, value := range events {
		if a.cache != nil && a.cache.Version > value.Version {
			continue
		}
		data = data + value.Data
		data = strings.ReplaceAll(data, "}{", ",")
	}
	var result map[string]interface{}
	e = json.Unmarshal([]byte(data), &result)
	if e != nil {
		return
	}
	bytes, e := json.Marshal(result)
	if e != nil {
		return
	}

	a.cache = &AggregateCache{
		UpdateTime: &events[len(events)].CreateTime,
		Data:       string(bytes),
		Version:    events[len(events)].Version,
	}
	return
}

/*
同步聚合和数据存储
*/
func (a *Aggregate) SyncCache() error {
	var t *time.Time
	var d string
	if a.cache == nil || a.cache.UpdateTime.IsZero() {
		t = nil
		d = ""
	} else {
		t = a.cache.UpdateTime
		d = a.cache.Data
	}
	//查找最后一个快照
	snap, e := snapshotStore.FindLastOrderByCreateTimeDesc(a.AggId, a.AggType, t)
	if e != nil {
		return e
	}
	if snap != nil && !snap.CreateTime.IsZero() {
		//快照数据为最新数据
		t = snap.CreateTime
		d = snap.Data
	}
	//t 4种情况
	//1,cache nil && snap ==nil       t == nil
	//2,cache !=nil && snap == nil    t == a.cache.updateTime
	//3,cache !=nil && snap !=nil     t == snap.createTime
	//4,cache nil && snap != nil      t == snap.createTime

	//根据time加载事件列表
	events, e := eventRepository.FindByOrderByAsc(a.AggId, a.AggType, t)
	if e != nil {
		return e
	}
	//事件溯源
	e = a.applyEvents(d, events)
	if e != nil {
		return e
	}
	return nil
}
