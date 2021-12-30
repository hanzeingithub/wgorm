package build_sqls

import (
	"reflect"
	"strings"
)

type SqlColumn struct {
	Name        string
	SqlOperator string
	SqlExpr     string
	SqlField    string
	Kind        reflect.Kind
}

func ParseSqlColumn(field reflect.StructField) (*SqlColumn, error) {
	sqlField := strings.TrimSpace(field.Tag.Get("sql_field"))
	sqlOperator := strings.TrimSpace(field.Tag.Get("sql_operator"))
	if sqlOperator == "" {
		sqlOperator = "="
	}
	if err := checkOperator(field, sqlOperator); err != nil {
		return nil, err
	}

	return &SqlColumn{
		SqlField:    sqlField,
		SqlOperator: sqlOperator,
		Name:        field.Name,
		Kind:        field.Type.Kind()}, nil
}
