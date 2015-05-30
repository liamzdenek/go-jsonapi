package jsonapi;

import (
    //"strings"
);

func(a *API) EntryFindRecordByResourceAndId(r *Request) {
    a.CentralFindRouter(r,
        r.Params.ByName("resource"),
        r.Params.ByName("id"),
        []string{},
        //OutputTypeResources, "",
    );
}

func(a *API) CentralFindRouter(r *Request, resourcestr, idstr string, preroute []string/*, outputtype OutputType, linkname string*/) {
    /*
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
        wctx.Push(work);
        work = NewTaskSingleLinkResolver(wctx, work, pre);
        //ii = ii.GetChild(pre);
    }
    attacher := NewTaskAttachIncluded(wctx, work, ii, outputtype, linkname);
    replyer := NewTaskReplyer(attacher);
    wctx.Push(work, attacher, replyer);
    a.Logger.Printf("Main Waiting\n");
    replyer.Wait();
    */
}
