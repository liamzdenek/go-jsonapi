package jsonapi;

import("net/http";)

type WorkFindLinksByIderTyperResult struct {
    Links *OutputLinkageSet
    Included *[]Record
}

type WorkFindLinksByIderTyper struct {
    IderTyper IderTyper
    II *IncludeInstructions
    Output chan chan *WorkFindLinksByIderTyperResult
    Result *WorkFindLinksByIderTyperResult
}

func NewWorkFindLinksByIderTyper(idertyper IderTyper, ii *IncludeInstructions) *WorkFindLinksByIderTyper {
    return &WorkFindLinksByIderTyper{
        II: ii,
        IderTyper: idertyper,
        Output: make(chan chan *WorkFindLinksByIderTyperResult),
    }
}

func (w *WorkFindLinksByIderTyper) Work(a *API, r *http.Request) {
    linker := NewLinkerDefault(
        a,
        a.RM.GetResource(w.IderTyper.Type()),
        w.IderTyper,
        r,
        w.II,
    );
    included := []Record{}
    w.Result = &WorkFindLinksByIderTyperResult{
        Links: linker.Link(&included),
        Included: &included,
    }
}

func(w *WorkFindLinksByIderTyper) ResponseWorker(has_paniced bool) {
    go func() {
        for r := range w.Output {
            r <- w.Result;
        }
    }()
}

func (w *WorkFindLinksByIderTyper) Cleanup(a *API, r *http.Request) {
    close(w.Output);
}

func(w *WorkFindLinksByIderTyper) GetResult() *WorkFindLinksByIderTyperResult  {
    r := make(chan *WorkFindLinksByIderTyperResult);
    defer close(r);
    w.Output <- r;
    return <-r;
}
