package handler

import (
	"fmt"
	"regexp"
	"testing"
)

func TestHandler(t *testing.T) {
	variables := &selectAggregate{}
	compile, e := regexp.Compile(`(?i).*\s*select \* from (\w+) where id='(\w+)'$`)
	if e != nil {
		panic(e.Error())
	}
	variables.compile = compile

	result, e := variables.Handler(`select * from UserAggregate where id='xudslajsdfhsadhf'`)
	fmt.Println(result, e)
}
