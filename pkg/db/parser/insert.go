package parser

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
)

type InsertParseResult struct {
	Columns   []string
	Values    [][]interface{}
	ValueMaps map[string][]interface{}
}

func ParseInsert(insert *sqlparser.Insert) (result *InsertParseResult, err error) {
	result = &InsertParseResult{}
	result.Columns = make([]string, len(insert.Columns))
	result.ValueMaps = make(map[string][]interface{}, len(insert.Columns))
	for key, value := range insert.Columns {
		result.Columns[key] = value.Lowered()
	}
	values, ok := insert.Rows.(sqlparser.Values)
	if !ok {
		err = fmt.Errorf("解析sql错误")
		return
	}
	tuples := []sqlparser.ValTuple(values)
	result.Values = make([][]interface{}, len(tuples))
	for itemKey, item := range tuples {
		exprs := sqlparser.Exprs(item)
		resultVal := make([]interface{}, len(exprs))
		for key, value := range exprs {
			val, ok := value.(*sqlparser.SQLVal)
			if !ok {
				err = fmt.Errorf("只支持值,不支持其他表达式")
				return
			}
			resultVal[key] = val.Val
		}
		result.Values[itemKey] = resultVal
	}
	for k, v := range result.Columns {
		result.ValueMaps[v] = result.Values[k]
	}
	return
}
