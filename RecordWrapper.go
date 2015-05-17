package jsonapi;

type RecordWrapper struct {
    Ider Ider
    Type_ string
    Linker Linker
    Show bool
}

func NewRecordWrapper(i Ider, t string, l Linker, show bool) *RecordWrapper {
    if i == nil {
        panic("NewRecordWrapper must not be provided with an Ider == nil");
    }
    return &RecordWrapper{
        Ider: i,
        Type_: t,
        Linker: l,
        Show: show,
    }
}

func(w *RecordWrapper) Id() string {
    return GetId(w.Ider);
}

func(w *RecordWrapper) SetId(s string) error {
    panic("TODO: This");
}

func(w *RecordWrapper) Link(included *[]Record) *OutputLinkageSet {
    return w.Linker.Link(included);
}

func(w *RecordWrapper) Type() string {
    return w.Type_;
}

func(w *RecordWrapper) Include() bool {
    return w.Show;
}

func(w RecordWrapper) Denature() interface{} {
    return w.Ider;
}
