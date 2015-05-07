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
    for _, inc := range *o.Included {
        d := &OutputDatum{Datum:inc};
        d.Prepare(nil);
        res = append(res,d);
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
