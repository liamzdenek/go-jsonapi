package jsonapi;

import (
    "strings"
);
func(a *API) EntryFindDefault(r *Request) {
    a.CentralFindRouter(r,
        r.Params.ByName("resource"),
        "",
        []string{},
        OutputTypeResources, "",
    );
}

func(a *API) EntryFindRecordByResourceAndId(r *Request) {
    a.CentralFindRouter(r,
        r.Params.ByName("resource"),
        r.Params.ByName("id"),
        []string{},
        OutputTypeResources, "",
    );
}

func(a *API) EntryFindRelationshipsByResourceId(r *Request) {
    a.CentralFindRouter(r,
        r.Params.ByName("resource"),
        r.Params.ByName("id"),
        []string{r.Params.ByName("linkname")},
        OutputTypeResources, "",
    );
}

func(a *API) EntryFindRelationshipByNameAndResourceId(r *Request) {
    /*a.CentralFindRouter(r,
        r.Params.ByName("resource"),
        r.Params.ByName("id"),
        []string{},
        OutputTypeLinkages,
        r.Params.ByName("secondlinkname"),
    );*/
}

func(a *API) CentralFindRouter(r *Request, resourcestr, idstr string, preroute []string, outputtype OutputType, linkname string) {
    resource := a.GetResource(resourcestr);
    pf := &PreparedFuture{
        Future: resource.GetFuture(),
        Relationship: nil,
    }
    fl := NewFutureList(r);

    var ids []string;
    if len(idstr) > 0 {
        ids = strings.Split(idstr, ",");
    }
    req := NewFutureRequest(r, &FutureRequestKindFindByIds{
        Ids: ids,
    });
    fl.PushFuture(pf);
    fl.PushRequest(pf,req);
    fl.Build(pf, resource, true).PrimaryData = pf.Future;
    fl.Takeover();
    defer fl.Defer();

    /*
    ids := []string{};
    if len(idstr) != 0 {
        ids = strings.Split(idstr,",");
    }
    var work TaskResultRecords = NewTaskFindByIds(
        resourcestr,
        ids,
        ii,
        "",
        NewPaginator(r),
    );
    for _,pre := range preroute {
        r.Push(work);
        work = NewTaskSingleLinkResolver(work, pre);
    }
    attacher := NewTaskAttachIncluded(work, ii, outputtype, linkname);
    replyer := NewTaskReplyer(attacher);
    r.Push(work, attacher, replyer);
    r.API.Logger.Infof("Main Waiting\n");
    replyer.Wait();*/
}

func(a *API) EntryDelete(r *Request) {
    /*deleter := NewTaskDelete(r.Params.ByName("resource"), r.Params.ByName("id"));
    r.Push(deleter);
    deleter.Wait();
    */
}

func(a *API) EntryCreate(r *Request) {
    /*creater := NewTaskCreate(r.Params.ByName("resource"), r.Params.ByName("id"));
    r.Push(creater);
    creater.Wait();
    */
}

func(a *API) EntryUpdate(r *Request) {
    /*creater := NewTaskUpdate(r.Params.ByName("resource"), r.Params.ByName("id"));
    r.Push(creater);
    creater.Wait();
    */
}
