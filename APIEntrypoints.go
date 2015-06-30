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
    a.CentralFindRouter(r,
        r.Params.ByName("resource"),
        r.Params.ByName("id"),
        []string{},
        OutputTypeLinkages,
        r.Params.ByName("secondlinkname"),
    );
}

func(a *API) CentralFindRouter(r *Request, resourcestr, idstr string, preroute []string, outputtype OutputType, linkname string) {
    resource := a.GetResource(resourcestr);
    ef := NewExecutableFuture(r, resource.GetFuture());

    var ids []string;
    if len(idstr) > 0 {
        ids = strings.Split(idstr, ",");
    }
    req := NewFutureRequest(r, &FutureRequestKindFindByIds{
        Ids: ids,
    });

    /*
    for _,pre := range preroute {
        r.API.Logger.Debugf("Get rel: %s %s\n", resource.Name, pre);
        relationship := a.GetRelationship(resource.Name, pre);
        resource = a.GetResource(relationship.DstResourceName);
        pf = &PreparedFuture{
            Future: resource.GetFuture(),
            Relationship: relationship,
            Parents: []*PreparedFuture{pf},
        };
        fl.PushFuture(pf);
    }
    */

    output := ef.Build(resource)
    output.PrimaryData = ef.Future
    output.PrimaryDataType = outputtype;
    //output.RelationshipName = linkname
    defer ef.Defer();
    ef.Takeover(req);

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
