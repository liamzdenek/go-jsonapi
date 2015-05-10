package jsonapi;

import ("reflect";"strings";);

type Denaturer interface {
    Denature() interface{}
}

func DenatureObject(data interface{}) map[string]interface{} {
    for {
        if d, found := data.(Denaturer); found {
            data = d.Denature();
        } else {
            break;
        }
    }
    v := reflect.Indirect(reflect.ValueOf(data));
    t := v.Type();

    values := make(map[string]interface{}, t.NumField());

    for i := 0; i < t.NumField(); i++ {
        var f string;
        if f = t.Field(i).Tag.Get("marshal-json"); len(f) == 0 {
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
                if(IsZeroOfUnderlyingType(v.Field(i).Interface())) {
                    continue;
                }
            }
        }
        values[tag[0]] = v.Field(i).Interface();
    }

    return values;
}

func IsZeroOfUnderlyingType(x interface{}) bool {
    return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
