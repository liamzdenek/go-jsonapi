package jsonapi;

import ("reflect";"strings";);

type Naturer interface {
    Nature() interface{}
}

func NatureObject(data map[string]interface{}, res interface{}) {
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
        if f = t.Field(i).Tag.Get("nature-json"); len(f) == 0 {
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
        //fmt.Printf("Setting field %v = %v\n", t.Field(i).Name, data[tag[0]]);
        v.Field(i).Set(
            reflect.ValueOf(data[tag[0]]).Convert(t.Field(i).Type),
        );
        delete(data, tag[0]);
    }
}
