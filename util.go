package jsonapi;

import ("reflect";"encoding/json");

/**
Check() will call panic(e) if the error provided is non-nil
*/
func Check(e error) {
    if e != nil {
        panic(e);
    }
}

/**
Reply() is just an alias for panic() -- it is syntax sugar in a few places.
*/
func Reply(a interface{}) {
    panic(a);
}

/**
Catch() functions similarly to most other languages try-catch... If a panic is thrown within the provided lambda, it will be intercepted and returned as an argument. If no error occurs, the return is nil.
*/
func Catch(f func()) (r interface{}) {
    defer func() {
        r = recover();
    }();
    f();
    return
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

func ParseJSONHelper(v *Record, raw []byte, t reflect.Type) (*Record, error) {
    if v.Attributes == nil {
        v.Attributes = reflect.New(t).Interface();
    }
    err := json.Unmarshal(raw, v);
    if(err != nil) {
        return nil, err;
    }
    if(rp.Data.Relationships == nil) {
        fmt.Printf("GOT NO RELATIONSHIPS\n");
        rp.Data.Relationships = &OutputLinkageSet{};
    }
    return v.(Ider), rp.Data.Id, &rp.Data.Type, rp.Relationships(), nil;
}
