package build_sqls

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hanzeingithub/wgorm/utils"
)

type WhereBuilder struct {
}

func (w *WhereBuilder) BuildSQL(condition interface{}) (query string, args []interface{}, err error) {
	rt, rv, err := utils.GetTypeAndValue(condition)
	if err != nil {
		return "", nil, err
	}
	return w.ParseSqlQuery(rt, rv)
}

func (w *WhereBuilder) ParseSqlQuery(rt reflect.Type, rv reflect.Value) (query string, values []interface{}, err error) {
	values = []interface{}{}
	isFirst := true
	queryBuilder := new(strings.Builder)
	for i := 0; i < rt.NumField(); i++ {
		column, err := ParseSqlColumn(rt.Field(i))
		if err != nil {
			return "", nil, err
		}
		param := rv.FieldByName(column.Name)
		// 没有值不处理
		if param.Kind() == reflect.Ptr && param.IsNil() {
			continue
		}
		// 如果是slice没有数据也不处理
		if param.Kind() == reflect.Slice && (param.IsNil() || param.Len() == 0) {
			continue
		}
		if param.Kind() == reflect.Ptr {
			param = param.Elem()
		}
		operator, placeHolder := getOperatorType(column.SqlOperator)
		if !isFirst {
			queryBuilder.WriteString(" AND ")
		} else {
			isFirst = false
		}
		queryBuilder.WriteString(fmt.Sprintf("`%s` %s %s", column.SqlField, operator, placeHolder))
		values = append(values, param.Interface())
	}
	query = queryBuilder.String()
	return
}
