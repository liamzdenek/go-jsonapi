package jsonapi;

import ("net/http";"fmt");

type TaskSingleLinkResolver struct {
    Context *TaskContext
    Parent TaskResultIderTypers
    Linkname string
    Output chan chan *TaskFindByIdsResult
    Result *TaskFindByIdsResult
}

func NewTaskSingleLinkResolver(ctx *TaskContext, t TaskResultIderTypers, linkname string) *TaskSingleLinkResolver {
    return &TaskSingleLinkResolver{
        Context: ctx,
        Parent: t,
        Linkname: linkname,
        Output: make(chan chan *TaskFindByIdsResult),
    }
}

func(t *TaskSingleLinkResolver) Work(a *API, r *http.Request) {
    result := t.Parent.GetResult();
    ii := NewIncludeInstructionsEmpty();
    ii.Push([]string{t.Linkname});
    data := []IderTyper{};
    for _, res := range result.Result {
        work := NewWorkFindLinksByIderTyper(res,ii);
        t.Context.Push(work);
        fmt.Printf("WORKRES: %#v\n", work.GetResult().Included);
        for _, inc := range *work.GetResult().Included {
            data = append(data, inc);
        }
    }
    t.Result = &TaskFindByIdsResult{
        Result: data,
        IsSingle: false, // TODO: fix this to get this data from the relationship
    };
}

func(t *TaskSingleLinkResolver) ResponseWorker(has_paniced bool) {
    go func() {
        for req := range t.Output {
            req <- t.Result;
        }
    }()
}

func(t *TaskSingleLinkResolver) Cleanup(a *API, r *http.Request) {
    close(t.Output);
}

func(t *TaskSingleLinkResolver) GetResult() *TaskFindByIdsResult {
    r := make(chan *TaskFindByIdsResult);
    defer close(r);
    t.Output <- r;
    return <-r;
}
