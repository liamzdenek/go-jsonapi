package jsonapi;

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
);

func(a *API) EntryFindRecordByResourceAndId(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    a.CentralFindRouter(a,w,r,
        ps.ByName("resource"),
        ps.ByName("id"),
        []string{},
        OutputTypeResources, "",
    );
}

func(a *API) CentralFindRouter(a *API, w http.ResponseWriter, r *http.Request, resourcestr, idstr string, preroute []string, outputtype OutputType, linkname string) {
    ii := NewIncludeInstructionsFromRequest(r);
    wctx := NewTaskContext(a,r,w,a.GetNewSession());
    defer wctx.Cleanup();
    ids := strings.Split(idstr,",");
    if(len(idstr) == 0) {
        ids = []string{}
    }
    var work TaskResultRecords = NewTaskFindByIds(
        resourcestr,
        ids,
        ii,
        "",
        NewPaginator(r),
    );
    /*
    for _,pre := range preroute {
        wctx.Push(work);
        work = NewTaskSingleLinkResolver(wctx, work, pre);
        //ii = ii.GetChild(pre);
    }*/
    attacher := NewTaskAttachIncluded(wctx, work, ii, outputtype, linkname);
    replyer := NewTaskReplyer(attacher);
    wctx.Push(work, attacher, replyer);
    a.Logger.Printf("Main Waiting\n");
    replyer.Wait();
}
