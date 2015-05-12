package jsonapi;

import "reflect";

func Check(e error) {
    if e != nil {
        panic(e);
    }
}

// Reply is a syntax sugar panic button
func Reply(a interface{}) {
    panic(a);
}

func GetField(field string, i interface{}) interface{} {
    return reflect.Indirect(reflect.ValueOf(i)).FieldByName(field).Interface();
}

func SetField(field string, i interface{}, v interface{}) {
    reflect.Indirect(reflect.ValueOf(i)).FieldByName(field).Set(reflect.ValueOf(v));
}
