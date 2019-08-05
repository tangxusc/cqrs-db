package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"regexp"
)

func init() {
	variables := &selectTxIsolation{}
	compile, e := regexp.Compile(`(?i).*\s*SELECT @@session.tx_isolation$`)
	if e != nil {
		panic(e.Error())
	}
	variables.compile = compile
	db.Handlers = append(db.Handlers, variables)
}

type selectTxIsolation struct {
	compile *regexp.Regexp
}

func (s *selectTxIsolation) Match(query string) bool {
	return s.compile.MatchString(query)
}

func (s *selectTxIsolation) Handler(query string) (*mysql.Result, error) {
	return nil, nil
}
