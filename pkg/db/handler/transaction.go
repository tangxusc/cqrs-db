package handler

import (
	"cqrs-db/pkg/db"
	"github.com/siddontang/go-mysql/mysql"
	"regexp"
)

func init() {
	variables := &transaction{}
	compile, e := regexp.Compile(`(?i).*\s*(start transaction|begin|commit|rollback)$`)
	if e != nil {
		panic(e.Error())
	}
	variables.compile = compile
	db.Handlers = append(db.Handlers, variables)
}

type transaction struct {
	compile *regexp.Regexp
}

func (s *transaction) Match(query string) bool {
	return s.compile.MatchString(query)
}

func (s *transaction) Handler(query string) (*mysql.Result, error) {
	return nil, nil
}
