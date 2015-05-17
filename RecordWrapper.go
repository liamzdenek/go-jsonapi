package jsonapi;

type RecordWrapper struct {
    Ider Ider
    Type_ string
    Context *TaskContext
    Work *WorkFindLinksByRecord
    II *IncludeInstructions
    ViaLinkName string
}

func NewRecordWrapper(i Ider, t string, ctx *TaskContext, vln string, ii *IncludeInstructions) *RecordWrapper {
    if i == nil {
        panic("NewRecordWrapper must not be provided with an Ider == nil");
    }
    if ii == nil {
        panic("NewRecordWrapper must not be provided with IncludeInstructions == nil")
    }
    if ctx == nil {
        panic("NewRecordWrapper must not be provided with TaskContext == nil");
    }
    if vln == "" {
        panic("NewRecordWrapper must not be provided with ViaLinkName == \"\"");
    }
    res := &RecordWrapper{
        Ider: i,
        Type_: t,
        Context: ctx,
        ViaLinkName: vln,
        II: ii,
    };
    work := NewWorkFindLinksByRecord(res,ii.GetChild(res.ViaLinkName));
    res.Context.Push(work);
    res.Work = work;
    return res;
}

func(w *RecordWrapper) Id() string {
    return GetId(w.Ider);
}

func(w *RecordWrapper) SetId(s string) error {
    return SetId(w.Ider, s);
}

func(w *RecordWrapper) Data() *WorkFindLinksByRecordResult {
    return w.Work.GetResult();
}

func(w *RecordWrapper) Type() string {
    return w.Type_;
}

func(w *RecordWrapper) Include() bool {
    return w.ViaLinkName == "" || w.II.ShouldInclude(w.ViaLinkName); // TODO chain this into shouldInclude w.II
}

func(w RecordWrapper) Denature() interface{} {
    return w.Ider;
}
