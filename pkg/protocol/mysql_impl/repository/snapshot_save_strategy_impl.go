package repository

import "github.com/tangxusc/cqrs-db/pkg/core"

type CountStrategy struct {
	Max int
}

func NewCountStrategy(max int) *CountStrategy {
	return &CountStrategy{Max: max}
}

func (s *CountStrategy) Allow(aggId string, aggType string, data string, events core.Events) bool {
	if len(events) > s.Max {
		return true
	}
	return false
}
