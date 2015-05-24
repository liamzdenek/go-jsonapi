package jsonapi;

import (
    "net/http";
    "strings";
    "github.com/julienschmidt/httprouter";
    //"errors";
    //"io/ioutil"
);

type RequestResolver struct{}

func NewRequestResolver() *RequestResolver {
    return &RequestResolver{};
}


/************************************************
 *
 * 
 */
func(rr *RequestResolver) HandlerRoot(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    defer func() {
        if raw := recover(); raw != nil {
            a.CatchResponses(w,r,raw);
        }
    }();
    o := NewOutput(r,a.Meta);
    Reply(o);
}
/************************************************
 *
 * HandlerFindResourceById is the entrypoint for /:resource requests, primarily:
 * * /user
 */
func(rr *RequestResolver) HandlerFindDefault(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    rr.CentralSearchRouter(a,w,r,
        ps.ByName("resource"),
        "",
        []string{},
        OutputTypeResources, "",
    );
}

/************************************************
 *
 * HandlerFindResourceById is the entrypoint for /:resource/:id requests, primarily:
 * * /user/1
 * * /user/1,2,3
 */

func(rr *RequestResolver) HandlerFindResourceById(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    rr.CentralSearchRouter(a,w,r,
        ps.ByName("resource"),
        ps.ByName("id"),
        []string{},
        OutputTypeResources, "",
    );
}


/************************************************
 *
 * HandlerFindLinksByResourceId  is the entrypoint for /:resource/:id/:linkname requests, primarily:
 * * /user/1/posts
 *
 * this handler does not support requests for multiple IDs (maybe it could?):
 * * /user/1,2,3/posts
 */

func(rr *RequestResolver) HandlerFindLinksByResourceId(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    rr.CentralSearchRouter(a,w,r,
        ps.ByName("resource"),
        ps.ByName("id"),
        []string{ps.ByName("linkname")},
        OutputTypeResources, "",
    );
}

/************************************************
 *
 * HandlerFindLinkByLinkNameAndResourceId is the entrypoint for
 * /:resource/:id/links/:secondlinkname requests, primarily:
 * * /user/1/links/posts
 *
 * this handler does not support requests for multiple IDs:
 * * /user/1,2,3/links/posts
 *
 * requests with :linkname != "links" will 404
 */

func(rr *RequestResolver) HandlerFindLinkByNameAndResourceId(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    rr.CentralSearchRouter(a,w,r,
        ps.ByName("resource"),
        ps.ByName("id"),
        []string{},
        OutputTypeLinkages,
        ps.ByName("secondlinkname"),
    );
}

func(rr *RequestResolver) CentralSearchRouter(a *API, w http.ResponseWriter, r *http.Request, resourcestr, idstr string, preroute []string, outputtype OutputType, linkname string) {
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
        "root",
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
}

/************************************************
 *
 * HandlerDelete is the entrypoint for
 * * DELETE /:resource/:id requests:
 *
 * this handler does not support requests for multiple IDs -- TODO: but it should:
 * * DELETE /user/1,2,3
 */

 // TODO: needs a diferent response when it does not exist per spec
func(rr *RequestResolver) HandlerDelete(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    wctx := NewTaskContext(a,r,w,a.GetNewSession());
    defer wctx.Cleanup();
    deleter := NewTaskDelete(
        ps.ByName("resource"),
        ps.ByName("id"),
    );
    wctx.Push(deleter);
    deleter.Wait();
}
/************************************************
 *
 * HandlerCreate is the entrypoint for:
 * * POST /:resource
 */
func(rr *RequestResolver) HandlerCreate(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    wctx := NewTaskContext(a,r,w,a.GetNewSession());
    defer wctx.Cleanup();
    creater := NewTaskCreate(
        ps.ByName("resource"),
        ps.ByName("id"),
    );
    wctx.Push(creater);
    creater.Wait();

    return;
}

func(rr *RequestResolver) HandlerUpdate(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
