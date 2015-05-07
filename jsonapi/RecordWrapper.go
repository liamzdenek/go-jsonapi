package jsonapi;

type RecordWrapper struct {
    Ider Ider
    Type_ string
    Linker Linker
}

func NewRecordWrapper(i Ider, t string, l Linker) *RecordWrapper {
    return &RecordWrapper{
        Ider: i,
        Type_: t,
        Linker: l,
    }
}

func(w *RecordWrapper) Id() string {
    return w.Ider.Id();
}

func(w *RecordWrapper) Link(included *[]Record) *OutputLinkageSet {
    return w.Linker.Link(included);
}

func(w *RecordWrapper) Type() string {
    return w.Type_;
}

func(w RecordWrapper) Denature() interface{} {
    return w.Ider;
}
