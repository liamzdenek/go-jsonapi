package jsonapi;

import ("reflect";"encoding/json";"fmt");

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

func ParseJSONHelper(v *Record, raw []byte, t reflect.Type) (*Record) {
    if v == nil {
        v = &Record{}
    }
    if v.Attributes == nil {
        v.Attributes = reflect.New(t).Interface();
    }
    rp := &RecordParserSimple{
        Data: v,
    };
    fmt.Printf("UNMARSHAL: %s\n",raw);
    err := json.Unmarshal(raw, rp);
    fmt.Printf("UNMARSHAL: %#v\n",rp);
    fmt.Printf("UNMARSHAL: %#v\n",rp.Data);
    fmt.Printf("UNMARSHAL: %#v\n",rp.Data.Relationships);
    fmt.Printf("UNMARSHAL: %#v\n",rp.Data.Attributes);
    if(err != nil) {
        panic(err);
    }
    return v;
}

func GetRelationshipDifferences(src, dst []OResourceIdentifier) (add, remove []OResourceIdentifier) {
    OUTER: for _,srcrid := range src {
        for _,dstrid := range dst {
            if srcrid.Id == dstrid.Id && srcrid.Type == dstrid.Type {
                continue OUTER;
            }
        }
        // if we got to this point, the link exists in the current set but does not exist in the new set, and must be deleted
        remove = append(remove, srcrid);
        panic("REMOVE");
    }
    
    // add ones that should be there now
    OUTER2: for _,dstrid := range dst {
        for _,srcrid := range src {
            if srcrid.Id == dstrid.Id && srcrid.Type == dstrid.Type {
                continue OUTER2;
            }
        }
        // if we got to this point, the link exists in the new set but does not exist in the old set, and must be added
        add = append(add, dstrid);
        panic("ADD");
    }
    return
}
