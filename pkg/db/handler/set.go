package handler

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/tangxusc/cqrs-db/pkg/db"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
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
	_, _, err := proxy.Proxy(query)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
