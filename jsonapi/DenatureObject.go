package jsonapi;

import ("reflect";"strings";);

func DenatureObject(data interface{}) map[string]interface{} {
    v := reflect.Indirect(reflect.ValueOf(data));
    t := v.Type();

    values := make(map[string]interface{}, t.NumField());

    for i := 0; i < t.NumField(); i++ {
        tag := strings.Split(t.Field(i).Tag.Get("json"), ",");
        if len(tag[0]) == 0 { 
            tag[0] = t.Field(i).Name
        }
        if len(tag) > 1 && len(tag[1]) > 0 {
            if(tag[1] == "omitempty") {
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