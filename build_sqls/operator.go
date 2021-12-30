package build_sqls

import (
	"fmt"
	"reflect"
	"strings"
)

var invalidOperators = map[string]bool{
	"<":        true,
	"<=":       true,
	"=":        true,
	"!=":       true,
	">":        true,
	">=":       true,
	"null":     true,
	"not null": true,
	"in":       true,
	"not in":   true,
}

func checkOperator(field reflect.StructField, sqlOperator string) error {
	if !invalidOperators[sqlOperator] {
		return fmt.Errorf("sql field %s sqlOperator %s invalid", field.Name, sqlOperator)
	}

	if field.Type.Kind() != reflect.Ptr {
		if field.Type.Kind() == reflect.Slice {
			return judgeOperatorArray(sqlOperator)
		} else {
			return fmt.Errorf("struct field %s should be pointer or slice", field.Name)
		}
	}
	return nil
}

func judgeOperatorArray(sqlOperator string) error {
	if isOperatorSupportArray(sqlOperator) {
		return nil
	}
	return fmt.Errorf("sqlOperator %s not support array", sqlOperator)
}

func isOperatorSupportArray(s string) bool {
	return strings.Contains(s, "in")
}

func getOperatorType(sqlOperator string) (operator string, placeHolder string) {
	if isOperatorSupportArray(sqlOperator) {
		return sqlOperator, "(?)"
	}
	if strings.Contains(sqlOperator, "null") {
		return fmt.Sprintf("is %s", sqlOperator), ""
	}
	return fmt.Sprintf("%s", sqlOperator), "?"
}
