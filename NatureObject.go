package jsonapi;

import ("reflect";"strings";"fmt";"errors");

type Naturer interface {
    Nature() interface{}
}

func NatureObject(data map[string]interface{}, res interface{}) error {
    //TODO: this function should catch panics
    for {
        if d, found := res.(Naturer); found {
            res = d.Nature();
        } else {
            break;
        }
    }
    v := reflect.Indirect(reflect.ValueOf(res));
    t := v.Type();

    //values := make(map[string]interface{}, t.NumField());
    fieldcount := t.NumField();
    for i := 0; i < fieldcount; i++ {
        // Never parse an ID into the object... that should remain in data
        if(strings.ToLower(t.Field(i).Name) == "id") {
            continue;
        }
        var f string;
        if f = t.Field(i).Tag.Get("unmarshal-json"); len(f) == 0 {
            f = t.Field(i).Tag.Get("json");
        }
        tag := strings.Split(f, ",");
        if len(tag[0]) == 0 {
            tag[0] = t.Field(i).Name
        }
        if(tag[0] == "-") {
            continue;
        }
        if len(tag) > 1 {
            if(len(tag[1]) > 0 && tag[1] == "omitempty") {
                if(IsZeroOfUnderlyingType(reflect.ValueOf(data[tag[0]]))) {
                    continue;
                }
            }
        }
        val := reflect.ValueOf(data[tag[0]]);
        target_type := t.Field(i).Type;
        if(!val.IsValid()) {
            return errors.New(fmt.Sprintf("Value received for field '%s' is not valid... did you forget to provide it?", tag[0]));
        }
        if(!val.Type().ConvertibleTo(target_type)) {
            return errors.New(fmt.Sprintf("Value retrieved for field '%s' is not ConvertibleTo Type '%s'", tag[0], target_type.String()));
        }
        v.Field(i).Set(val.Convert(target_type));
        delete(data, tag[0]);
    }
    return nil;
}
