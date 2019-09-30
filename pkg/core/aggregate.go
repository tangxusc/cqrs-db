package core

import (
	"context"
	"fmt"
)

/*
聚合
*/
type Aggregate struct {
	Version      int
	eventsChan   chan Events
	recoveryChan chan Events
	ctx          context.Context
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
	select {
	case <-a.ctx.Done():
		return fmt.Errorf("已关闭")
	default:
		select {
		case a.eventsChan <- events:
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
