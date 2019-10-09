package memory

import (
	"context"
	"fmt"
	"github.com/ReneKroon/ttlcache"
	"github.com/tangxusc/cqrs-db/pkg/core"
	"time"
)

type AggregateManagerImpl struct {
	ctx       context.Context
	container *ttlcache.Cache
}

func NewAggregateManagerImpl(ctx context.Context) *AggregateManagerImpl {
	aggm := &AggregateManagerImpl{ctx: ctx}
	cache := ttlcache.NewCache()
	cache.SetTTL(time.Minute * 10)
	aggm.container = cache
	return aggm
}

type AggregateEntry struct {
	agg    *core.Aggregate
	ctx    context.Context
	cancel context.CancelFunc
}

func NewAggregateEntry(ctx context.Context) *AggregateEntry {
	entry := &AggregateEntry{}
	ctx2, cancel := context.WithCancel(ctx)
	entry.ctx = ctx2
	entry.cancel = cancel
	return entry
}

func (a *AggregateManagerImpl) Get(aggId, aggType string) *core.Aggregate {
	target, exist := a.container.Get(key(aggId, aggType))
	if exist {
		return target.(*AggregateEntry).agg
	}
	entry := NewAggregateEntry(a.ctx)
	defer func() {
		if e := recover(); e != nil {
			entry.cancel()
		}
	}()
	aggregate, e := core.NewAggregate(aggId, aggType, entry.ctx)
	if e != nil {
		panic(e)
	}
	entry.agg = aggregate
	a.container.Set(key(aggId, aggType), entry)
	return aggregate
}

func key(aggId string, aggType string) string {
	return fmt.Sprintf(`%s-%s`, aggType, aggId)
}
