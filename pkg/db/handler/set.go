package handler

import (
	"cqrs-db/pkg/db"
	"github.com/siddontang/go-mysql/mysql"
	"regexp"
)

func init() {
	variables := &set{}
	compile, e := regexp.Compile(`(?i).*\s*set .*$`)
	if e != nil {
		panic(e.Error())
	}
	variables.compile = compile
	db.Handlers = append(db.Handlers, variables)
}

type set struct {
	compile *regexp.Regexp
}

func (s *set) Match(query string) bool {
	return s.compile.MatchString(query)
}

func (s *set) Handler(query string) (*mysql.Result, error) {
	return nil, nil
}
