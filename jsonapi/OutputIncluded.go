package jsonapi;

import("encoding/json";"fmt";);

type OutputIncluded struct {
    Included *[]Record
}

func NewOutputIncluded(included *[]Record) *OutputIncluded {
    fmt.Printf("NewOutputIncluded: %#v\n", included);
    return &OutputIncluded{
        Included: included,
    }
}

func(o OutputIncluded) MarshalJSON() ([]byte, error) {
    res := []interface{}{};
    todo_list := *o.Included
    var inc Record
    for {
        if(len(todo_list) >= 1) {
            inc = todo_list[0]
        } else {
            break;
        }
        fmt.Printf("DATUM: %#v\n", DenatureObject(inc));
        d := &OutputDatum{Datum:inc};
        d.Prepare(&todo_list);
        if(inc.Include()) {
            fmt.Printf("DATUM INCLUDED\n");
            res = append(res,d);
        }
        if(len(todo_list) >= 1) {
            todo_list = todo_list[1:];
        }
    }
    return json.Marshal(res);
    //return json.Marshal(o.Included);
}
/*
func(o *OutputIncluded) Push(included ...Record) {
    for _,ider := range included {
        o.Included = append(o.Included, ider);
    }
}*/

func(o *OutputIncluded) ShouldBeVisible() bool {
    // TODO: the spec requires more complicated visibility logic than this
    fmt.Printf("Should be visible: %s\n", o.Included);
    return len(*o.Included) > 0
}
