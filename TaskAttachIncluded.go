package jsonapi;

import("net/http";"fmt");

type TaskAttachIncluded struct {
    Context *TaskContext
    Parent TaskResultIderTypers
    II *IncludeInstructions
    Output chan chan *Output
    ActualOutput *Output
}

func NewTaskAttachIncluded(ctx *TaskContext, parent TaskResultIderTypers, ii *IncludeInstructions) *TaskAttachIncluded {
    return &TaskAttachIncluded{
        Context: ctx,
        Parent: parent,
        II: ii,
        Output: make(chan chan *Output),
    }
}

func (w *TaskAttachIncluded) Work(a *API, r *http.Request) {
    result := w.Parent.GetResult();
    queue := result.Result;
    data := []*OutputDatum{};
    included := []Record{};
    first := true
    for {
        tqueue := queue;
        queue = []IderTyper{}
        d := map[IderTyper]*WorkFindLinksByIderTyper{};
        for _, idertyper := range tqueue {
            work := NewWorkFindLinksByIderTyper(idertyper,w.II);
            w.Context.Push(work);
            d[idertyper] = work;
        }
        for idertyper, work := range d {
            result := work.GetResult();
            fmt.Printf("RESULT: %#v\n",result);
            if(first) {
                data = append(data, &OutputDatum{
                    Datum: NewRecordWrapper(idertyper, idertyper.Type(),NewLinkerStatic(result.Links), true),
                });
            }
            for _, record := range *result.Included {
                queue = append(queue, record);
                fmt.Printf("GOT RECORD: %#v\n", record);
                if(record.Include()) {
                    included = append(included, record)
                }
            }
        }
        first = false;
        if len(queue) == 0 {
            break;
        }
    }
    res := &Output{};
    fmt.Printf("DATA: %#v\n", data);
    res.Data = NewOutputDataResources(result.IsSingle, data);
    res.Included = NewOutputIncluded(&included);

    w.ActualOutput = res;
}

func (w *TaskAttachIncluded) ResponseWorker(has_paniced bool) {
    go func() {
        for req := range w.Output {
            req <- w.ActualOutput;
        }
    }();
}

func (w *TaskAttachIncluded) Cleanup(a *API, r *http.Request) {
    fmt.Printf("TaskAttachIncluded.Cleanup\n");
    close(w.Output);
}

func (w *TaskAttachIncluded) GetResult() *Output {
    r := make(chan *Output);
    defer close(r);
    w.Output <- r;
    return <-r;
}
