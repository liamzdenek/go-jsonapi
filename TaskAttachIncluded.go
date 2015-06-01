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

func (t *TaskAttachIncluded) Work(r *Request) {
    parent_result := t.Parent.GetResult();
    r.API.Logger.Debugf("PARENT: %#v\n", parent_result);
    res := NewOutput();

    output_primary := []*Record{};
    output_included := []*Record{};

    queue := parent_result.Records;
    primary_data_count := len(queue);
    for {
        if(len(queue) == 0) {
            break;
        }
        var next *Record;
        next, queue = queue[0], queue[1:]; // queue pop

        r.API.Logger.Infof("MAIN LOOP HANDLING: %#v\n", next);
        //relationships := next.GetRelationships();
        if(primary_data_count > 0) {
            primary_data_count--;
            output_primary = append(output_primary, next);
        } else if(next.ShouldInclude) {
            output_included = append(output_included, next);
        }
        rels := next.GetRelationships()
        queue = append(queue, rels.Included...)
        next.Relationships = rels.Relationships;
    }

    if(t.OutputType == OutputTypeResources) {
        res.Data = &ORecords{
            IsSingle: parent_result.IsSingle,
            Records: output_primary,
        };
    } else {
        /*
        res.Data = &ORelationship{
            IsSingle: parent_result.IsSingle, // TODO: check this, i don't think it's right?
            Data: ConvertToResourceIdentifiers(parent
        }
        */
        panic("TODO");
    }
    res.Included = output_included;
    t.ActualOutput = res;
    /*
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
    */
}

func (w *TaskAttachIncluded) ResponseWorker(has_paniced bool) {
    go func() {
        for req := range w.Output {
            req <- w.ActualOutput;
        }
    }();
}

func (w *TaskAttachIncluded) Cleanup(r *Request) {
    r.API.Logger.Debugf("TaskAttachIncluded.Cleanup\n");
    close(w.Output);
}

func (w *TaskAttachIncluded) GetResult() *Output {
    r := make(chan *Output);
    defer close(r);
    w.Output <- r;
    return <-r;
}
