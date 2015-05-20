package jsonapi;

import ("reflect";"encoding/json";);

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
    for {
        if n,ok := i.(Denaturer); ok {
            i = n.Denature();
        } else {
            break;
        }
    }
    return reflect.Indirect(reflect.ValueOf(i)).FieldByName(field).Interface();
}

func SetField(field string, i interface{}, v interface{}) {
    reflect.Indirect(reflect.ValueOf(i)).FieldByName(field).Set(reflect.ValueOf(v));
}

func ParseJSONHelper(raw []byte, t reflect.Type) (Ider, *string, *string, *OutputLinkageSet, error) {
    v := reflect.New(t).Interface();
    rp := NewRecordParserSimple(v);
    err := json.Unmarshal(raw, rp);
    if(err != nil) {
        return nil, nil, nil, nil, err;
    }
    if(rp.Data.Linkages == nil) {
        rp.Data.Linkages = &OutputLinkageSet{};
    }
    return v.(Ider), rp.Data.Id, &rp.Data.Type, rp.Linkages(), nil;
}

func Catch(f func()) (r interface{}) {
    defer func() {
        r = recover();
    }();
    f();
    return
}
