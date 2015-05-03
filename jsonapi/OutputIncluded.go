package jsonapi;

import("encoding/json");

type OutputIncluded struct {
    Included []IderTyper
}

func(o OutputIncluded) MarshalJSON() ([]byte, error) {
    return json.Marshal(o.Included);
}

func(o *OutputIncluded) Push(included ...IderTyper) {
    for _,ider := range included {
        o.Included = append(o.Included, ider);
    }
}
