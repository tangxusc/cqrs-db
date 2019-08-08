package aggregate

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var Locks = &sync.Map{}

type LockEntry struct {
	Value int32
	Key   string
}

func (l *LockEntry) Lock() {
	for {
		swapped := atomic.CompareAndSwapInt32(&l.Value, 0, 1)
		if swapped {
			return
		}
	}
}

func (l *LockEntry) Unlock() {
	for {
		swapped := atomic.CompareAndSwapInt32(&l.Value, 1, 0)
		if swapped {
			return
		}
	}
}

func (l *LockEntry) delete() {
	if atomic.LoadInt32(&l.Value) == 0 {
		Locks.Delete(l.Key)
	}
}

/*
解锁聚合
*/
func UnLock(id string, aggType string) {
	UnLockWithKey(getKey(id, aggType))
}

func UnLockWithKey(key string) {
	value, ok := Locks.Load(key)
	if !ok {
		return
	}
	entry := value.(*LockEntry)
	defer func() {
		entry.delete()
	}()
	entry.Unlock()
}

/*
lock聚合,防止并发
*/
func Lock(id string, aggType string) string {
	key := getKey(id, aggType)
	entry := &LockEntry{
		Value: 0,
		Key:   key,
	}
	existLock, loaded := Locks.LoadOrStore(key, entry)
	if loaded {
		entry = existLock.(*LockEntry)
	}
	entry.Lock()
	return key
}

func getKey(id string, aggType string) string {
	return fmt.Sprintf("%s:%s", aggType, id)
}
