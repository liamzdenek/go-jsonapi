package jsonapi;

import("net/http";"fmt")

type WorkFindLinksByRecordResult struct {
    Links *OutputLinkageSet
    Included *[]Record
}

type WorkFindLinksByRecord struct {
    Record Record
    II *IncludeInstructions
    Output chan chan *WorkFindLinksByRecordResult
    Result *WorkFindLinksByRecordResult
}

func NewWorkFindLinksByRecord(idertyper Record, ii *IncludeInstructions) *WorkFindLinksByRecord {
    return &WorkFindLinksByRecord{
        II: ii,
        Record: idertyper,
        Output: make(chan chan *WorkFindLinksByRecordResult),
    }
}

func (w *WorkFindLinksByRecord) Work(a *API, r *http.Request) {
    fmt.Printf("GOT RECORD TO FIND LINKS: %#v\n", w.Record);
    linker := NewLinkerDefault(
        a,
        a.RM.GetResource(w.Record.Type()),
        w.Record,
        r,
        w.II,
    );
    included := []Record{}
    w.Result = &WorkFindLinksByRecordResult{
        Links: linker.Link(&included),
        Included: &included,
    }
}

func(w *WorkFindLinksByRecord) ResponseWorker(has_paniced bool) {
    go func() {
        for r := range w.Output {
            r <- w.Result;
        }
    }()
}

func (w *WorkFindLinksByRecord) Cleanup(a *API, r *http.Request) {
    close(w.Output);
}

func(w *WorkFindLinksByRecord) GetResult() *WorkFindLinksByRecordResult  {
    r := make(chan *WorkFindLinksByRecordResult);
    defer close(r);
    w.Output <- r;
    return <-r;
}
