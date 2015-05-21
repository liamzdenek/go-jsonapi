package jsonapi;

import ("net/http";);

type TaskSingleLinkResolver struct {
    Context *TaskContext
    Parent TaskResultRecords
    Linkname string
    Output chan chan *TaskResultRecordData
    Result *TaskResultRecordData
}

func NewTaskSingleLinkResolver(ctx *TaskContext, t TaskResultRecords, linkname string) *TaskSingleLinkResolver {
    return &TaskSingleLinkResolver{
        Context: ctx,
        Parent: t,
        Linkname: linkname,
        Output: make(chan chan *TaskResultRecordData),
    }
}

// TODO: make the parent_name in this function passed as an arg to NewTaskSingleLinkResolver instead of determining it from the result, as the current setup could create inconsistent behavior and is inherently incompatible with multi-type resources
func(t *TaskSingleLinkResolver) Work(a *API, s Session, ctx *TaskContext, r *http.Request) {
    result := t.Parent.GetResult();
    ii := NewIncludeInstructionsEmpty();
    ii.Push([]string{t.Linkname});
    data := []Record{};
    parent_name := "";
    for _, res := range result.Result {
        parent_name = res.Type();
        work := NewWorkFindLinksByRecord(res,ii);
        t.Context.Push(work);
        a.Logger.Printf("WORKRES: %#v\n", work.GetResult().Included);
        for _, inc := range *work.GetResult().Included {
            data = append(data, inc);
        }
    }
    isSingle := false;
    rel := a.RM.GetRelationship(parent_name, t.Linkname);
    if(rel != nil) {
        isSingle = rel.B.IsSingle()
    }
    t.Result = &TaskResultRecordData{
        Result: data,
        IsSingle: isSingle,
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

func(t *TaskSingleLinkResolver) GetResult() *TaskResultRecordData {
    r := make(chan *TaskResultRecordData);
    defer close(r);
    t.Output <- r;
    return <-r;
}