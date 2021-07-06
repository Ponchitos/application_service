package common

import "reflect"

func IsPointer(value interface{}) bool {
	rValue := reflect.ValueOf(value)

	return rValue.Kind() == reflect.Ptr
}
