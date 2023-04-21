package refutil

import (
	"fmt"
	"reflect"
)

func ValueOf(v interface{}) reflect.Value {
	if val, ok := v.(reflect.Value); ok {
		// check if it's already a reflect.Value.
		return val
	}

	return reflect.ValueOf(v)
}

func IsFunc(kindAble interface{ Kind() reflect.Kind }) bool {
	return kindAble.Kind() == reflect.Func
}

func IsBlank(data interface{}) bool {
	value := reflect.ValueOf(data)
	fmt.Println(value.IsValid())
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
