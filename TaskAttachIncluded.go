package jsonapi;

import("net/http";"fmt");

type TaskAttachIncluded struct {
    Context *TaskContext
    Parent TaskResultRecords
    II *IncludeInstructions
    Output chan chan *Output
    ActualOutput *Output
    OutputType OutputType
    Linkname string
}


type OutputType int;

const (
    OutputTypeResources OutputType = iota;
    OutputTypeLinkages
);

func NewTaskAttachIncluded(ctx *TaskContext, parent TaskResultRecords, ii *IncludeInstructions, outputtype OutputType, linkname string) *TaskAttachIncluded {
    return &TaskAttachIncluded{
        Context: ctx,
        Parent: parent,
        II: ii,
        Output: make(chan chan *Output),
        OutputType: outputtype,
        Linkname: linkname,
    }
}

func (w *TaskAttachIncluded) Work(a *API, r *http.Request) {
    result := w.Parent.GetResult();
    queue := result.Result;
    data := []*OutputDatum{};
    linkage := OutputLinkage{};
    included := []Record{};
    first := true
    for {
        tqueue := queue;
        queue = []Record{}
        d := map[Record]*WorkFindLinksByRecord{};
        for _, idertyper := range tqueue {
            work := NewWorkFindLinksByRecord(idertyper,w.II);
            w.Context.Push(work);
            d[idertyper] = work;
        }
        for idertyper, work := range d {
            result := work.GetResult();
            if(first) {
                if w.OutputType == OutputTypeResources {
                    data = append(data, &OutputDatum{
                        Datum: NewRecordWrapper(idertyper, idertyper.Type(),NewLinkerStatic(result.Links), true),
                    });
                } else {
                    for _, links := range result.Links.Linkages {
                        if(links.LinkName == w.Linkname) {
                            for _, link := range links.Links {
                                linkage.Links = append(linkage.Links, link);
                            }
                        }
                    }
                }
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
    if w.OutputType == OutputTypeResources {
        res.Data = NewOutputDataResources(result.IsSingle, data);
    } else {
        fmt.Printf("\nLINKAGE: %#v\n\n", linkage);
        res.Data = NewOutputDataLinkage(result.IsSingle, &linkage);
    }
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
