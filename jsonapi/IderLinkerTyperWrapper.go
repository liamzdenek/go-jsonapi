package jsonapi;

type IderLinkerTyperWrapper struct {
    Ider Ider
    Type_ string
    Linker Linker
}

func NewIderLinkerTyperWrapper(i Ider, t string, l Linker) *IderLinkerTyperWrapper {
    return &IderLinkerTyperWrapper{
        Ider: i,
        Type_: t,
        Linker: l,
    }
}

func(w *IderLinkerTyperWrapper) Id() string {
    return w.Ider.Id();
}

func(w *IderLinkerTyperWrapper) Link() *OutputLinkageSet {
    return w.Linker.Link();
}

func(w *IderLinkerTyperWrapper) Type() string {
    return w.Type_;
}

func(w IderLinkerTyperWrapper) Denature() interface{} {
    return w.Ider;
}
