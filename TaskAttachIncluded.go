package jsonapi;

type TaskAttachIncluded struct {
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

func NewTaskAttachIncluded(parent TaskResultRecords, ii *IncludeInstructions, outputtype OutputType, linkname string) *TaskAttachIncluded {
    return &TaskAttachIncluded{
        Parent: parent,
        II: ii,
        Output: make(chan chan *Output),
        OutputType: outputtype,
        Linkname: linkname,
    }
}

func (w *TaskAttachIncluded) Work(r *Request) {
    parent := w.Parent.GetResult();
    output_primary := []Record{};
    output_included := []Record{};
    output_linkage := OutputLinkage{};
    queue := []Record{};
    queue = append(queue, parent.Result...);
    primary_count := len(queue);
    for {
        if len(queue) == 0 {
            break;
        }
        var next Record;
        next, queue = queue[0], queue[1:]; // queue pop
        result := next.Data();
        
        if(primary_count > 0) {
            primary_count--;
            a.Logger.Printf("Pushing To Primary: %#v\n", next);
            output_primary = append(output_primary, next);
            for _, links := range next.Data().Links.Linkages {
                if(links.LinkName == w.Linkname) {
                    for _, link := range links.Links {
                        output_linkage.Links = append(output_linkage.Links, link);
                    }
                }
            }
        } else if(next.Include()) {
            a.Logger.Printf("Pushing To Included: %#v\n", next);
            output_included = append(output_included, next);
        }

        for _,included := range *result.Included {
            a.Logger.Printf("Pushing To Queue: %#v\n", included);
            queue = append(queue, included);
        }
    }

    res := NewOutput(nil);
    if w.OutputType == OutputTypeResources {
        a.Logger.Printf("Primary data is a resource");
        res.Data = NewOutputDataResources(parent.IsSingle, output_primary);
    } else {
        a.Logger.Printf("Primary data is a linkage");
        res.Data = NewOutputDataLinkage(parent.IsSingle, &output_linkage);
    }

    a.Logger.Printf("PAGINATOR: %#v\n", parent.Paginator);

    res.SetPaginator(r,parent.Paginator);

    res.Included.Included = &output_included;
    w.ActualOutput = res;
}

func (w *TaskAttachIncluded) ResponseWorker(has_paniced bool) {
    go func() {
        for req := range w.Output {
            req <- w.ActualOutput;
        }
    }();
}

func (w *TaskAttachIncluded) Cleanup(r *Request) {
    r.API.Logger.Printf("TaskAttachIncluded.Cleanup\n");
    close(w.Output);
}

func (w *TaskAttachIncluded) GetResult() *Output {
    r := make(chan *Output);
    defer close(r);
    w.Output <- r;
    return <-r;
}
