package repository

import "github.com/tangxusc/cqrs-db/pkg/core"

type CountStrategy struct {
	current uint
	Max     uint
}

func NewCountStrategy(max uint) *CountStrategy {
	return &CountStrategy{Max: max}
}

func (s *CountStrategy) Allow(aggId string, aggType string, data map[string]interface{}, events core.Events) bool {
	s.current = uint(len(events)) + s.current
	if s.current > s.Max {
		return true
	}
	return false
}

type CountStrategyFactory struct {
	Max uint
}

func NewCountStrategyFactory(max uint) *CountStrategyFactory {
	return &CountStrategyFactory{Max: max}
}

func (s *CountStrategyFactory) NewStrategyInstance() core.SnapshotSaveStrategy {
	return NewCountStrategy(s.Max)
}
