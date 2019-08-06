package aggregate

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var locks = &sync.Map{}

type lockEntry struct {
	value int32
}

func (l *lockEntry) Lock() {
	for {
		swapped := atomic.CompareAndSwapInt32(&l.value, 0, 1)
		if swapped {
			return
		}
	}
}

func (l *lockEntry) Unlock() {
	for {
		swapped := atomic.CompareAndSwapInt32(&l.value, 1, 0)
		if swapped {
			return
		}
	}
}

/*
解锁聚合
*/
func UnLock(id string, aggType string) {
	defer locks.Delete(getKey(id, aggType))
	value, ok := locks.Load(getKey(id, aggType))
	if !ok {
		return
	}
	entry := value.(*lockEntry)
	entry.Unlock()
}

/*
lock聚合,防止并发
*/
func Lock(id string, aggType string) {
	entry := &lockEntry{
		value: 0,
	}
	existLock, loaded := locks.LoadOrStore(getKey(id, aggType), entry)
	if loaded {
		entry = existLock.(*lockEntry)
	}
	entry.Lock()
}

func getKey(id string, aggType string) string {
	return fmt.Sprintf("%s:%s", aggType, id)
}
