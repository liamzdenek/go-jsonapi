package jsonapi

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Ider interface {
	Id() string
	SetId(string) error
}

func GetId(ider interface{}) string {
	if ider == nil {
		panic("IDER provided to GetId CANNOT be nil")
	}
	if manual, ok := ider.(Ider); ok {
		return manual.Id()
	}
	field, _ := GetIdField(ider)
	val := field.Interface()
	if str, ok := val.(string); ok {
		return str
	}
	if id, ok := val.(int); ok {
		return fmt.Sprintf("%d", id)
	}
	if str, ok := val.(fmt.Stringer); ok {
		return str.String()
	}
	panic("Couldn't properly format string")
}

func SetId(ider interface{}, id string) error {
	f, _ := GetIdField(ider)
	//t := f.Type()
	if _, ok := f.Interface().(string); ok {
		f.Set(reflect.ValueOf(id))
		return nil
	}
	return errors.New("SetId does not have a mapping for converting between these types")
}

func GetIdField(ider interface{}) (reflect.Value, reflect.StructField) {
	return GetFieldByTag(ider, "id")
}

func GetFieldByTag(ider interface{}, realtag string) (reflect.Value, reflect.StructField) {
	var val reflect.Value
	var typ reflect.Type
	for {
		val = reflect.Indirect(reflect.ValueOf(ider))
		typ = val.Type()
		ider = val.Interface()
		if val.Kind() != reflect.Ptr {
			break
		}
	}
	fields := val.NumField()
	for i := 0; i < fields; i++ {
		tags := strings.Split(typ.Field(i).Tag.Get("jsonapi"), ",")
		for _, tag := range tags {
			if tag == realtag {
				return val.Field(i), typ.Field(i)
			}
		}
	}
	panic(fmt.Sprintf("Couldn't get field \"%s\" for provided ider: %#v\n", realtag, ider))
}
