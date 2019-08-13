package event

import (
	"fmt"
	"github.com/tangxusc/cqrs-db/pkg/db/parser"
	"github.com/tangxusc/cqrs-db/pkg/proxy"
	"github.com/tangxusc/cqrs-db/pkg/util"
	"strings"
)

var Columns = []string{"id", "type", "agg_id", "agg_type", "create_time", "data"}

func SaveEvent(result *parser.InsertParseResult) error {
	if len(result.Values) == 0 {
		return fmt.Errorf("值不能为空")
	}
	if len(result.Columns) == 0 {
		result.Columns = Columns
	} else {
		result.Columns = append(result.Columns, "id")
	}
	columnsSql := strings.Join(result.Columns, ",")
	for key, vValue := range result.Values {
		result.Values[key] = append(vValue, util.GenerateUuid())
	}
	return proxy.Inserts(`insert into event(`+columnsSql+`) values (?,?,?,?,?,?)`, result.Values)
}
