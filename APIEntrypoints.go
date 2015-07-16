package jsonapi;

import (
    "strings"
);
func(a *API) EntryFindDefault(r *Request) {
    a.CentralFindRouter(r,
        r.Params.ByName("resource"),
        "",
        []string{},
        OutputTypeResources,
    );
}

func(a *API) EntryFindRecordByResourceAndId(r *Request) {
    a.CentralFindRouter(r,
        r.Params.ByName("resource"),
        r.Params.ByName("id"),
        []string{},
        OutputTypeResources,
    );
}

func(a *API) EntryFindRelationshipsByResourceId(r *Request) {
    a.CentralFindRouter(r,
        r.Params.ByName("resource"),
        r.Params.ByName("id"),
        []string{r.Params.ByName("linkname")},
        OutputTypeResources,
    );
}

func(a *API) EntryFindRelationshipByNameAndResourceId(r *Request) {
    a.CentralFindRouter(r,
        r.Params.ByName("resource"),
        r.Params.ByName("id"),
        []string{r.Params.ByName("secondlinkname")},
        OutputTypeLinkages,
    );
}

func(a *API) CentralFindRouter(r *Request, resourcestr, idstr string, preroute []string, outputtype OutputType) {
    resource := a.GetResource(resourcestr);
    ef := NewExecutableFuture(r, resource.GetFuture());
    ef.Resource = resource;

    first_ef := ef;
    var ids []string;
    if len(idstr) > 0 {
        ids = strings.Split(idstr, ",");
    }
    req := NewFutureRequest(r, &FutureRequestKindFindByIds{
        Ids: ids,
    });

    for _,pre := range preroute {
        r.API.Logger.Debugf("Get rel: %s %s\n", resource.Name, pre);
        relationship := a.GetRelationship(resource.Name, pre);
        resource = a.GetResource(relationship.DstResourceName);
        nef := NewExecutableFuture(r,resource.GetFuture());
        nef.Resource = resource
        nef.Relationship = relationship
        ef.PushChild(nef.Relationship, nef);
        ef = nef;
    }

    output := ef.Build(resource)
    output.PrimaryData = ef.Future
    output.PrimaryDataType = outputtype;
    defer first_ef.Defer();
    first_ef.Takeover(req);
}

func(a *API) EntryDelete(r *Request) {
    idstr := r.Params.ByName("id");
    var ids []string;
    if len(idstr) > 0 {
        ids = strings.Split(idstr, ",");
    }

    resource := a.GetResource(r.Params.ByName("resource"));
    ef := NewExecutableFuture(r, resource.GetFuture());
    ef.Takeover(NewFutureRequest(r, &FutureRequestKindDeleteByIds{
        Ids: ids,
    }));
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
