package build_sqls

import (
	"fmt"
	"reflect"

	"gouqinggan/wgorm/utils"
)

type UpdateBuilder struct {
}

func (w *UpdateBuilder) BuildSQL(condition interface{}) (updatesMap map[string]interface{}, err error) {
	rt, rv, err := utils.GetTypeAndValue(condition)
	if err != nil {
		return nil, err
	}
	return w.ParseSqlQuery(rt, rv)
}

func (w *UpdateBuilder) ParseSqlQuery(rt reflect.Type, rv reflect.Value) (updatesMap map[string]interface{}, err error) {
	updatesMap = make(map[string]interface{})
	for i := 0; i < rt.NumField(); i++ {
		column, err := ParseSqlColumn(rt.Field(i))
		if err != nil {
			return nil, err
		}
		param := rv.FieldByName(column.Name)
		if err = judgeUpdateParam(param); err != nil {
			return nil, err
		}
		// 没有值不处理
		if param.IsNil() {
			continue
		}
		param = param.Elem()
		updatesMap[column.SqlField] = param.Interface()
	}
	return
}

func judgeUpdateParam(param reflect.Value) (err error) {
	if param.Kind() != reflect.Ptr {
		return fmt.Errorf("update param must be ptr")
	}
	if param.Elem().Kind() == reflect.Slice {
		return fmt.Errorf("update param can't be slice")
	}
	return
}
