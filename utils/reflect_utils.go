package utils

import (
	"fmt"
	"reflect"
)

// GetTypeAndValue target 需要是指针，并且必须是struct
func GetTypeAndValue(target interface{}) (reflectType reflect.Type, reflectValue reflect.Value, err error) {
	reflectValue = reflect.ValueOf(target)
	if !reflectValue.IsValid() {
		return nil, reflectValue, fmt.Errorf("where condition invaild")
	}
	if reflectValue.Kind() != reflect.Ptr {
		return nil, reflectValue, fmt.Errorf("where condition should be ptr")
	}
	reflectValue = reflectValue.Elem()
	if reflectValue.Kind() != reflect.Struct {
		fmt.Println(reflectValue.Kind())
		return nil, reflectValue, fmt.Errorf("where condition should be struct")
	}
	reflectType = reflectValue.Type()
	return
}
