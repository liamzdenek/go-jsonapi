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

func (w *TaskAttachIncluded) Work(ctx *TaskContext, a *API, r *http.Request) {
    parent := w.Parent.GetResult();
    output_primary := []Record{};
    output_included := []Record{};
    queue := []Record{};
    queue = append(queue, parent.Result...);
    primary_count := len(queue);
    for {
        if len(queue) == 0 {
            break;
        }
        var next Record;
        next, queue = queue[0], queue[1:];
        result := next.Data();
        
        if(primary_count > 0) {
            primary_count--;
            a.Logger.Printf("Pushing To Primary: %#v\n", next);
            output_primary = append(output_primary, next);
        } else if(next.Include()) {
            a.Logger.Printf("Pushing To Included: %#v\n", next);
            output_included = append(output_included, next);
        }

        for _,included := range *result.Included {
            a.Logger.Printf("Pushing To Queue: %#v\n", included);
            queue = append(queue, included);
        }
    }

    res := NewOutput(r);
    res.Data = NewOutputDataResources(parent.IsSingle, output_primary);
    res.Included.Included = &output_included;
    w.ActualOutput = res;

    //panic("TODO: FIX DIS");
    /*
    queue := result.Result;
    data := []*OutputDatum{};
    linkage := OutputLinkage{};
    included := []Record{};
    first := true
    for {

        tqueue := queue;
        queue = []Record{}
        d := map[Record]*WorkFindLinksByRecord{};
        for _, record := range tqueue {
            work := NewWorkFindLinksByRecord(record,w.II);
            w.Context.Push(work);
            d[record] = work;
        }
        for record, work := range d {
            result := work.GetResult();
            a.Logger.Printf("ATTACH INCLUDED GOT RESULT: %#v %#v %s %s\n\n", result.Links, result.Included, record.Type(), GetId(record));
            if(first) {
                //if w.OutputType == OutputTypeResources {
                    data = append(data, &OutputDatum{
                        Datum: record,
                    });
                //} else {
                //    for _, links := range result.Links.Linkages {
                //        if(links.LinkName == w.Linkname) {
                //            for _, link := range links.Links {
                //                linkage.Links = append(linkage.Links, link);
                //            }
                //        }
                //    }
                //}
            }
            for _, crecord := range *result.Included {
                a.Logger.Printf("INCLUDING RECORD: %#v\n", crecord);
                queue = append(queue, crecord);
                if(crecord.Include()) {
                    included = append(included, crecord)
                }
            }
        }
        first = false;
        if len(queue) == 0 {
            break;
        }
    }
    res := &Output{};
    fmt.Printf("ACTUAL OUTPUT: %#v\n", data);
    if w.OutputType == OutputTypeResources {
        res.Data = NewOutputDataResources(result.IsSingle, data);
    } else {
        res.Data = NewOutputDataLinkage(result.IsSingle, &linkage);
    }
    res.Included = NewOutputIncluded(&included);
    
    w.ActualOutput = res;
    */
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
