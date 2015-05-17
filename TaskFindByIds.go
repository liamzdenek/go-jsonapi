package jsonapi;

import("fmt";"net/http")

type TaskFindByIds struct {
    Resource string
    Ids []string
    Output chan chan *TaskResultRecordData
    Result *TaskResultRecordData
    II *IncludeInstructions
    ViaLinkName string
}

func NewTaskFindByIds(resource string, ids []string, ii *IncludeInstructions, vln string) *TaskFindByIds {
    return &TaskFindByIds{
        Output: make(chan chan *TaskResultRecordData),
        Ids: ids,
        Resource: resource,
        II: ii,
        ViaLinkName: vln,
    }
}

func(w *TaskFindByIds) Work(wctx *TaskContext, a *API, r *http.Request) {
    resource := a.RM.GetResource(w.Resource);

    if(resource == nil) {
        panic(NewResponderErrorResourceDoesNotExist(w.Resource));
    }

    // TODO: make this a loop over all the IDs
    for _, id := range w.Ids {
        resource.A.Authenticate("resource.FindOne."+w.Resource, id, r);
    }

    data := []Ider{}

    var err error;
    if(len(w.Ids) == 0) {
        panic("Ids must be longer than 0");
    } else if(len(w.Ids) == 1) {
        var ider Ider;
        ider, err = resource.R.FindOne(w.Ids[0]);
        if ider != nil {
            data = []Ider{ider}
        }
    } else {
        data, err = resource.R.FindMany(w.Ids);
    }
    if err != nil {
        // TODO: is this the right error?
        panic(NewResponderError(err));
    }
    fmt.Printf("GOT DATA: %#v\n", data);
    res := []Record{};
    for _,ider := range data {
        res = append(res, NewRecordWrapper(
            ider,
            w.Resource,
            wctx,
            w.ViaLinkName,
            w.II,
        ));
    }
    w.Result = &TaskResultRecordData{
        Result: res,
        IsSingle: len(w.Ids) == 1,
    }
}

func(w *TaskFindByIds) ResponseWorker(has_paniced bool) {
    go func() {
        for out := range w.Output {
            out <- w.Result;
        }
    }();
}

func(w *TaskFindByIds) Cleanup(a *API, r *http.Request) {
    fmt.Printf("TASKFINDBYIDS CLEANUP\n");
    close(w.Output);
}

func(w *TaskFindByIds) GetResult() *TaskResultRecordData {
    r := make(chan *TaskResultRecordData);
    defer close(r);
    w.Output <- r;
    return <-r;
}
