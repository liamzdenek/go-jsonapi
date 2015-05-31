package jsonapi;

import (
    "strings"
);

func(a *API) EntryFindRecordByResourceAndId(r *Request) {
    a.CentralFindRouter(r,
        r.Params.ByName("resource"),
        r.Params.ByName("id"),
        []string{},
        OutputTypeResources, "",
    );
}

func(a *API) CentralFindRouter(r *Request, resourcestr, idstr string, preroute []string, outputtype OutputType, linkname string) {
    ii := r.IncludeInstructions;
    ids := strings.Split(idstr,",");
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
    replyer.Wait();
}
