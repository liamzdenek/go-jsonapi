package jsonapi;

type IderLinkerTyperWrapper struct {
    Ider Ider
    Type_ string
}

func NewIderLinkerTyperWrapper(i Ider, t string) *IderLinkerTyperWrapper {
    return &IderLinkerTyperWrapper{
        Ider: i,
        Type_: t,
    }
}

func(w *IderLinkerTyperWrapper) Id() string {
    return w.Ider.Id();
}

func(w *IderLinkerTyperWrapper) Link() *OutputLinkageSet {
    //panic("TODO");
    return nil;
}

func(w *IderLinkerTyperWrapper) Type() string {
    return w.Type_;
}

func(w IderLinkerTyperWrapper) Denature() interface{} {
    return w.Ider;
}
