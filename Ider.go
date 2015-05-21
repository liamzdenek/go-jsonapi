package jsonapi;

import ("reflect";"strings";"fmt";"errors")

type Ider interface{

}

func GetId(ider Ider) string {
    if ider == nil {
        panic("IDER provided to GetId CANNOT be nil");
    }
    if manual, ok := ider.(IderManual); ok {
        return manual.Id();
    }
    val := GetIdField(ider).Interface();
    if str, ok := val.(string); ok {
        return str;
    }
    if id, ok := val.(int); ok {
        return fmt.Sprintf("%d", id);
    }
    if str, ok := val.(fmt.Stringer); ok {
        return str.String();
    }
    panic("Couldn't properly format string");
}

func SetId(ider Ider, id string) error {
    f := GetIdField(ider);
    //t := f.Type()
    if _, ok := f.Interface().(string); ok {
        f.Set(reflect.ValueOf(id));
        return nil;
    }
    return errors.New("SetId does not have a mapping for converting between these types");
}

func GetIdField(ider Ider) (reflect.Value) {
    return GetFieldByTag(ider, "id");
}

func GetFieldByTag(ider Ider, realtag string) (reflect.Value) {
    val := reflect.Indirect(reflect.ValueOf(ider))
    typ := val.Type();
    fields := val.NumField();
    for i := 0; i < fields; i++ {
        tags := strings.Split(typ.Field(i).Tag.Get("jsonapi"),",");
        for _,tag := range tags {
            if(tag == realtag) {
                return val.Field(i);
            }
        }
    }
    panic(fmt.Sprintf("Couldn't get field \"%s\" for provided ider: %#v\n", realtag, ider));
}
