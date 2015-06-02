package jsonapi;

type TaskFindByIds struct {
    Resource string
    Ids []string
    Output chan chan *TaskResultRecordData
    Result *TaskResultRecordData
    II *IncludeInstructions
    ViaLinkName string
    Paginator Paginator
}

func NewTaskFindByIds(resource string, ids []string, ii *IncludeInstructions, vln string, Paginator Paginator) *TaskFindByIds {
    return &TaskFindByIds{
        Output: make(chan chan *TaskResultRecordData),
        Ids: ids,
        Resource: resource,
        II: ii,
        ViaLinkName: vln,
        Paginator: Paginator,
    }
}

func(t *TaskFindByIds) Work(r *Request) {
    resource := r.API.GetResource(t.Resource);

    if(resource == nil) {
        panic(NewResponderErrorResourceDoesNotExist(t.Resource));
    }

    // TODO: make this a loop over all the IDs
    for _, id := range t.Ids {
        resource.Authenticator.Authenticate(r,"resource.FindOne."+t.Resource, id);
    }

    data := []*Record{}
    rp := RequestParams{
        Paginator: t.Paginator,
    };

    var err error;
    if(len(t.Ids) == 0) {
        data,err = resource.Resource.FindDefault(r,rp)
    } else if(len(t.Ids) == 1) {
        var record *Record;
        record, err = resource.Resource.FindOne(r,rp,t.Ids[0]);
        if record != nil {
            data = []*Record{record}
        }
    } else {
        data, err = resource.Resource.FindMany(r,rp, t.Ids);
    }
    if err != nil {
        // TODO: is this the right error?
        panic(NewResponderBaseErrors(500, err));
    }
    r.API.Logger.Debugf("GOT DATA: %#v\n", data);
    for _,record := range data {
        record.PrepareRelationships(r, t.II.GetChild(t.ViaLinkName));
    }
    t.Result = &TaskResultRecordData{
        Records: data,
        Paginator: &t.Paginator,
        IsSingle: len(t.Ids) == 1,
    }
}

func(t *TaskFindByIds) ResponseWorker(has_paniced bool) {
    go func() {
        for out := range t.Output {
            out <- t.Result;
        }
    }();
}

func(t *TaskFindByIds) Cleanup(r *Request) {
    r.API.Logger.Debugf("TASKFINDBYIDS CLEANUP\n");
    close(t.Output);
}

func(t *TaskFindByIds) GetResult() *TaskResultRecordData {
    r := make(chan *TaskResultRecordData);
    defer close(r);
    t.Output <- r;
    return <-r;
}
