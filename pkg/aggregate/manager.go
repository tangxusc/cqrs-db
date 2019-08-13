package aggregate

import (
	"fmt"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"sync"
	"time"
)

var SourceMap = sync.Map{}

func GetSource(id, aggType string, handler *db.ConnHandler) (*Source, error) {
	if !handler.TxBegin {
		return nil, fmt.Errorf("需要开启事务才能查询聚合")
	}
	if len(handler.TxKey) != 0 {
		return nil, fmt.Errorf("一个事务中只能查询一次聚合")
	}
	key := getKey(id, aggType)
	source := &Source{
		Id:       id,
		AggType:  aggType,
		Key:      key,
		locked:   0,
		lifeTime: time.Minute * 10,
	}
	actual, loaded := SourceMap.LoadOrStore(key, source)
	if loaded {
		return actual.(*Source), nil
	}
	source.start()
	return source, nil
}

func getKey(id string, aggType string) string {
	return fmt.Sprintf("%s:%s", aggType, id)
}

func GetSourceByKey(key string) {
	value, ok := SourceMap.Load(key)
	if ok {
		value.(*Source).unlock()
	}
}
