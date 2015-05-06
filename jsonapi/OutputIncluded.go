package jsonapi;

import("encoding/json";"fmt");

type OutputIncluded struct {
    Included *[]IderTyper
}

func NewOutputIncluded(included *[]IderTyper) *OutputIncluded {
    fmt.Printf("NewOutputIncluded: %#v\n", included);
    return &OutputIncluded{
        Included: included,
    }
}

func(o OutputIncluded) MarshalJSON() ([]byte, error) {
    return json.Marshal(o.Included);
}
/*
func(o *OutputIncluded) Push(included ...IderTyper) {
    for _,ider := range included {
        o.Included = append(o.Included, ider);
    }
}*/

func(o *OutputIncluded) ShouldBeVisible() bool {
    // TODO: the spec requires more complicated visibility logic than this
    fmt.Printf("Should be visible: %s\n", o.Included);
    return len(*o.Included) > 0
}
