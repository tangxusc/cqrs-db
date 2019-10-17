package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type AggregateCache struct {
	UpdateTime *time.Time
	Data       map[string]interface{}
	Version    int
}

func NewAggregateCache() *AggregateCache {
	return &AggregateCache{
		UpdateTime: nil,
		Data:       make(map[string]interface{}),
		Version:    0,
	}
}

/*
聚合
*/
type Aggregate struct {
	AggId                string
	AggType              string
	eventsChan           chan Events
	recoveryChan         chan Events
	ctx                  context.Context
	snapshotSaveStrategy SnapshotSaveStrategy

	cache *AggregateCache
}

func NewAggregate(aggId, aggType string, ctx context.Context) (*Aggregate, error) {
	agg := &Aggregate{
		AggId:                aggId,
		AggType:              aggType,
		eventsChan:           make(chan Events),
		recoveryChan:         make(chan Events),
		cache:                NewAggregateCache(),
		snapshotSaveStrategy: snapshotSaveStrategyFactory.NewStrategyInstance(),
	}
	e := agg.SyncCache()
	if e != nil {
		return nil, e
	}
	go agg.Start(ctx)

	return agg, nil
}

func (a *Aggregate) Start(ctx context.Context) {
	a.ctx = ctx
	//panic时,也可以关闭chan
	defer func() {
		close(a.eventsChan)
		close(a.recoveryChan)
		if e := recover(); e != nil {
			logrus.Errorf("[Aggregate] Aggregate[%s-%s] error:%v", a.AggType, a.AggId, e)
		}
	}()
	for {
		select {
		//优先级1
		case <-ctx.Done():
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
					//读取无缓存chan,处理完一个才读取下一个,不存在并发
					//计算快照,并缓存data
					e := a.applyEvents(a.cache.Data, events)
					if e != nil {
						panic(e)
					}
					go func() {
						allow := a.snapshotSaveStrategy.Allow(a.AggId, a.AggType, a.cache.Data, events)
						if allow {
							snapshotStore.Save(a.AggId, a.AggType, a.cache)
						}
					}()
					go a.sendEvent(events)
				default:
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
		if e != nil {
			logrus.Warnf("[Aggregate] %s-%s sendEvent error:%v , begin retry...", a.AggType, a.AggId, e)
		}
	}
}

func (a *Aggregate) PutSendChan(events Events) error {
	minVersion := minVersion(events)
	if a.cache.Version >= minVersion {
		return fmt.Errorf(`[Aggregate]聚合[%s-%s]版本错误,当前版本:%v,传入的最小版本:%v`, a.AggType, a.AggId, a.cache.Version, minVersion)
	}
	//保存event
	err := eventRepository.SaveEvents(a, events)
	if err != nil {
		return err
	}
	select {
	case <-a.ctx.Done():
		return fmt.Errorf("已关闭")
	default:
		select {
		case a.eventsChan <- events:
			return nil
		}
	}
}

func minVersion(events Events) int {
	var min = events[0].Version
	for _, v := range events {
		if v.Version < min {
			min = v.Version
		}
	}
	return min
}

func (a *Aggregate) PutRecoveryChan(events Events) error {
	select {
	case <-a.ctx.Done():
		return fmt.Errorf("已关闭")
	default:
		select {
		case a.recoveryChan <- events:
			return nil
		}
	}
}

/*
1,查询快照
2,根据快照查询事件
3,进行json merge溯源
4,返回结果
*/
func (a *Aggregate) Sourcing() (map[string]interface{}, int, error) {
	return a.cache.Data, a.cache.Version, nil
}

func Sourcing(id, aggType string) (data map[string]interface{}, version int, e error) {
	agg, e := aggregateCache.Get(id, aggType)
	if e != nil {
		return nil, -1, e
	}
	if agg == nil {
		return nil, 0, nil
	}
	return agg.Sourcing()
}

func (a *Aggregate) applyEvents(data map[string]interface{}, events []*Event) (e error) {
	if len(events) == 0 {
		return
	}
	var aggData string
	if len(data) > 0 {
		marshal, e := json.Marshal(data)
		if e != nil {
			return e
		}
		aggData = string(marshal)
	} else {
		aggData = ""
	}
	//按照顺序合并数据
	//聚合版本字段处理,如果小于当前版本,则需要丢弃这个event
	for _, value := range events {
		if a.cache != nil && a.cache.Version > value.Version {
			continue
		}
		aggData = aggData + value.Data
		aggData = strings.ReplaceAll(aggData, "}\n{", ",")
		aggData = strings.ReplaceAll(aggData, "}{", ",")
	}
	var result map[string]interface{}
	e = json.Unmarshal([]byte(aggData), &result)
	if e != nil {
		return
	}

	a.cache = &AggregateCache{
		UpdateTime: &events[len(events)-1].CreateTime,
		Data:       result,
		Version:    maxVersion(events),
	}
	return
}

func maxVersion(events []*Event) int {
	var max int
	for _, v := range events {
		if v.Version > max {
			max = v.Version
		}
	}
	return max
}

/*
同步聚合和数据存储
*/
func (a *Aggregate) SyncCache() error {
	var t *time.Time
	var d map[string]interface{}
	if a.cache == nil || a.cache.UpdateTime == nil || a.cache.UpdateTime.IsZero() {
		t = nil
		d = make(map[string]interface{})
	} else {
		t = a.cache.UpdateTime
		d = a.cache.Data
	}
	//查找最后一个快照
	snap, e := snapshotStore.FindLastOrderByCreateTimeDesc(a.AggId, a.AggType, t)
	if e != nil {
		logrus.Errorf("[aggregate]查找快照出现错误,将使用全量event溯源,错误详情:%v", e)
	}
	if e != nil && snap != nil && snap.CreateTime != nil && !snap.CreateTime.IsZero() {
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
