package jsonapi;

import ("encoding/json");

type IderTyperWrapper struct {
    Ider Ider
    Type_ string
}

func NewIderTyperWrapper(i Ider, t string) *IderTyperWrapper {
    return &IderTyperWrapper{
        Ider: i,
        Type_: t,
    }
}

func(w *IderTyperWrapper) Id() string {
    return w.Ider.Id();
}

func(w *IderTyperWrapper) Type() string {
    return w.Type_;
}

func(w IderTyperWrapper) Denature() interface{} {
    return w.Ider;
}

func(w IderTyperWrapper) MarshalJSON() ([]byte, error) {
    res := DenatureObject(w);
    delete(res, "ID");
    delete(res, "Id");
    delete(res, "iD");
    res["id"] = w.Id();
    res["type"] = w.Type();
    return json.Marshal(res);
}
