package jsonapi;

type IderTyperWrapper struct {
    Ider Ider
    Typestr string
}

func NewIderTyperWrapper(ider Ider, typestr string) *IderTyperWrapper {
    return &IderTyperWrapper{
        Ider: ider,
        Typestr: typestr,
    }
}

func(w *IderTyperWrapper) Id() string {
    return GetId(w.Ider);
}

func(w *IderTyperWrapper) SetId(s string) error {
    panic("TODO: this");
}

func(w *IderTyperWrapper) Type() string {
    return w.Typestr;
}

func(w *IderTyperWrapper) Denature() interface{} {
    return w.Ider;
}
