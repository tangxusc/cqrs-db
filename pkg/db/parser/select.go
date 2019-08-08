package parser

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
)

type SelectParseResult struct {
	TableName   string
	TableAsName string
	//名称,别名
	ColumnMap map[string]string
	Where     []string
}

func ParseSelect(stmt *sqlparser.Select) (result *SelectParseResult) {
	result = &SelectParseResult{}
	err := ParseTableName(stmt, result)
	if err != nil {
		return
	}
	err = ParseColumnName(stmt, result)
	if err != nil {
		return
	}
	err = ParseWhere(stmt, result)
	if err != nil {
		return
	}
	return
}

func ParseWhere(stmt *sqlparser.Select, result *SelectParseResult) error {
	if stmt.Where == nil {
		return nil
	}
	expr, ok := stmt.Where.Expr.(*sqlparser.ComparisonExpr)
	if !ok {
		return nil
	}
	val := expr.Right.(*sqlparser.SQLVal)
	result.Where = make([]string, 1)
	result.Where[0] = string(val.Val)
	return nil
}

func ParseColumnName(selectVar *sqlparser.Select, result *SelectParseResult) (err error) {
	exprs := []sqlparser.SelectExpr(selectVar.SelectExprs)
	result.ColumnMap = make(map[string]string, len(exprs))
	for _, value := range exprs {
		switch value.(type) {
		case *sqlparser.AliasedExpr:
			expr := value.(*sqlparser.AliasedExpr)
			colName, ok := expr.Expr.(*sqlparser.ColName)
			if !ok {
				err = fmt.Errorf("不支持的列名称")
				return
			}
			colNameString := colName.Name.Lowered()
			colAsNameString := colNameString
			if !expr.As.IsEmpty() {
				colAsNameString = expr.As.Lowered()
			}
			result.ColumnMap[colNameString] = colAsNameString
		case *sqlparser.StarExpr:
		//对于*如何处理
		default:
			err = fmt.Errorf("不支持nextval等函数")
			return
		}
	}
	return
}

func ParseTableName(selectVar *sqlparser.Select, result *SelectParseResult) (err error) {
	exprs := []sqlparser.TableExpr(selectVar.From)
	if len(exprs) != 1 {
		err = fmt.Errorf("不支持查询多个表")
		return
	}
	expr := exprs[0]
	switch expr.(type) {
	case *sqlparser.AliasedTableExpr:
		tableExpr := expr.(*sqlparser.AliasedTableExpr)
		result.TableAsName = tableExpr.As.String()
		name, ok := tableExpr.Expr.(sqlparser.TableName)
		if !ok {
			err = fmt.Errorf("不支持子查询等操作")
			return
		}
		result.TableName = name.Name.String()
	default:
		err = fmt.Errorf("不支持的sql查询")
		return
	}
	return
}
