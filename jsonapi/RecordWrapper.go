package jsonapi;

type RecordWrapper struct {
    Ider Ider
    Type_ string
    Linker Linker
    Show bool
}

func NewRecordWrapper(i Ider, t string, l Linker, show bool) *RecordWrapper {
    return &RecordWrapper{
        Ider: i,
        Type_: t,
        Linker: l,
        Show: show,
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

func(w *RecordWrapper) Include() bool {
    return w.Show;
}

func(w RecordWrapper) Denature() interface{} {
    return w.Ider;
}
