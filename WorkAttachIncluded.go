package jsonapi;

import("net/http";"fmt");

type WorkAttachIncluded struct {
    Context *WorkerContext
    Parent WorkerResultIderTypers
    II *IncludeInstructions
    Output chan chan *Output
    ActualOutput *Output
}

func NewWorkAttachIncluded(ctx *WorkerContext, parent WorkerResultIderTypers, ii *IncludeInstructions) *WorkAttachIncluded {
    return &WorkAttachIncluded{
        Context: ctx,
        Parent: parent,
        II: ii,
        Output: make(chan chan *Output),
    }
}

func (w *WorkAttachIncluded) Work(a *API, r *http.Request) {
    fmt.Printf("GETTING PARENT RESULT\n");
    todo := w.Parent.GetResult();
    fmt.Printf("PARENT: %#v\n", todo);
    data := []*OutputDatum{};
    included := []Record{};
    first := true
    for {
        ttodo := todo;
        todo = []IderTyper{}
        d := map[IderTyper]*WorkFindLinksByIderTyper{};
        for _, idertyper := range ttodo {
            work := NewWorkFindLinksByIderTyper(idertyper,w.II);
            PushWork(w.Context, work);
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
                todo = append(todo, record);
                fmt.Printf("GOT RECORD: %#v\n", record);
                if(record.Include()) {
                    included = append(included, record)
                }
            }
        }
        first = false;
        if len(todo) == 0 {
            break;
        }
    }
    res := &Output{};
    fmt.Printf("DATA: %#v\n", data);
    // TODO: repalce const true with actual calculation for single
    res.Data = NewOutputDataResources(false, data);
    res.Included = NewOutputIncluded(&included);

    w.ActualOutput = res;
}

func (w *WorkAttachIncluded) ResponseWorker(has_paniced bool) {
    go func() {
        for req := range w.Output {
            req <- w.ActualOutput;
        }
    }();
}

func (w *WorkAttachIncluded) Cleanup(a *API, r *http.Request) {
    fmt.Printf("WorkAttachIncluded.Cleanup\n");
    close(w.Output);
}

func (w *WorkAttachIncluded) GetResult() *Output {
    r := make(chan *Output);
    defer close(r);
    w.Output <- r;
    return <-r;
}
