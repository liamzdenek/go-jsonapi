package jsonapi;

type TaskFindLinksByRecordResult struct {
    Relationships *ORelationships
    Included []*Record
}

type TaskFindLinksByRecord struct {
    Record *Record
    II *IncludeInstructions
    Output chan chan *TaskFindLinksByRecordResult
    Result *TaskFindLinksByRecordResult
}

func NewTaskFindLinksByRecord(r *Record, ii *IncludeInstructions) *TaskFindLinksByRecord {
    return &TaskFindLinksByRecord{
        II: ii,
        Record: r,
        Output: make(chan chan *TaskFindLinksByRecordResult),
    }
}

func (t *TaskFindLinksByRecord) Work(r *Request) {
    //a.Logger.Printf("GOT RECORD TO FIND LINKS: %#v\n", w.Record.Link);
    //resource := r.API.GetResource(t.Record.Type);
    result := &TaskFindLinksByRecordResult{
        Relationships: &ORelationships{
            Relationships: []*ORelationship{},
            RelatedBase: r.GetBaseURL()+t.Record.Type+"/"+t.Record.Id,
        },
        Included: []*Record{},
    }
    for linkname,relationship := range r.API.GetRelationshipsByResource(t.Record.Type) {
        shouldFetch := t.II.ShouldFetch(linkname);
        or, included := relationship.Resolve(r, t.Record, shouldFetch, t.II);
        result.Relationships.Relationships = append(result.Relationships.Relationships, or);
        result.Included = append(result.Included, included...);
    }
    t.Result = result;
    /*
    linker := NewLinkerDefault(
        a,
        s,
        a.RM.GetResource(w.Record.Type()),
        w.Record,
        wctx,
        r,
        w.II,
    );
    included := &[]Record{}
    t.Result = &TaskFindLinksByRecordResult{
        Links: linker.Link(included),
        Included: included,
    }
    */
    //a.Logger.Printf("GOT RECORD LINKS: %#v\n", w.Result);
}

func(t *TaskFindLinksByRecord) ResponseWorker(has_paniced bool) {
    go func() {
        for r := range t.Output {
            r <- t.Result;
        }
    }()
}

func (t *TaskFindLinksByRecord) Cleanup(r *Request) {
    close(t.Output);
}

func(t *TaskFindLinksByRecord) GetResult() *TaskFindLinksByRecordResult  {
    r := make(chan *TaskFindLinksByRecordResult);
    defer close(r);
    t.Output <- r;
    res := <-r;
    return res;
}
