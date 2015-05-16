package jsonapi;

import ("net/http";"strings";"github.com/julienschmidt/httprouter";"fmt";"errors";"io/ioutil");

type RequestResolver struct{}

func NewRequestResolver() *RequestResolver {
    return &RequestResolver{};
}

/************************************************
 *
 * HandlerFindResourceById is the entrypoint for /:resource/:id requests, primarily:
 * * /user/1
 * * /user/1,2,3
 */

func(rr *RequestResolver) HandlerFindResourceById(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    ii := NewIncludeInstructionsFromRequest(r);
    wctx := NewTaskContext(a,r,w);
    defer wctx.Cleanup();
    work := NewTaskFindByIds(
        ps.ByName("resource"),
        strings.Split(ps.ByName("id"),","),
    );
    attacher := NewTaskAttachIncluded(wctx, work, ii, OutputTypeResources);
    replyer := NewTaskReplyer(attacher);
    wctx.Push(work, attacher, replyer);
    fmt.Printf("Main Waiting\n");
    replyer.Wait();
}


/************************************************
 *
 * HandlerFindLinksByResourceId  is the entrypoint for /:resource/:id/links requests, primarily:
 * * /user/1/links
 *
 * this handler does not support requests for multiple IDs (maybe it could?):
 * * /user/1,2,3/links
 */

func(rr *RequestResolver) HandlerFindLinksByResourceId(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    ii := NewIncludeInstructionsFromRequest(r);
    wctx := NewTaskContext(a,r,w);
    defer wctx.Cleanup();
    primary := NewTaskFindByIds(
        ps.ByName("resource"),
        strings.Split(ps.ByName("id"),","),
    );
    single := NewTaskSingleLinkResolver(wctx, primary, ps.ByName("linkname"));
    attacher := NewTaskAttachIncluded(wctx, single, ii, OutputTypeResources);
    replyer := NewTaskReplyer(attacher);
    wctx.Push(primary, single, attacher, replyer);
    fmt.Printf("Main Waiting\n");
    replyer.Wait();
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
 * requests with :linkname = "links" will 404
 */

func(rr *RequestResolver) HandlerFindLinkByNameAndResourceId(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    ii := NewIncludeInstructionsFromRequest(r);
    wctx := NewTaskContext(a,r,w);
    defer wctx.Cleanup();
    primary := NewTaskFindByIds(
        ps.ByName("resource"),
        strings.Split(ps.ByName("id"),","),
    );
    attacher := NewTaskAttachIncluded(wctx, primary, ii, OutputTypeLinkages);
    replyer := NewTaskReplyer(attacher);
    wctx.Push(primary, attacher, replyer);
    fmt.Printf("Main Waiting\n");
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

func(rr *RequestResolver) HandlerDelete(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    ids := strings.Split(ps.ByName("id"),",");
    isSingle := len(ids) == 1;
    if(!isSingle) {
        panic(NewResponderError(errors.New("This request does not support more than one id")));
    }

    resource_str := ps.ByName("resource");
    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(NewResponderErrorResourceDoesNotExist(resource_str));
    }

    resource.A.Authenticate("resource.Delete."+resource_str, ids[0], r);

    err := resource.R.Delete(ids[0]);
    Check(err);
    Reply(NewResponderResourceSuccessfullyDeleted());
}
/************************************************
 *
 * HandlerCreate is the entrypoint for:
 * * POST /:resource
 */
func(rr *RequestResolver) HandlerCreate(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    ctx := a.GetNewContext();
    resource_str := ps.ByName("resource");
    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(NewResponderErrorResourceDoesNotExist(resource_str));
    }

    resource.A.Authenticate("resource.Create."+resource_str, "", r);

    body, err := ioutil.ReadAll(r.Body);
    if err != nil {
        panic(NewResponderError(errors.New(fmt.Sprintf("Body could not be parsed: %v\n", err))));
    }

    ider,id,rtype,linkages,err := resource.R.ParseJSON(body);
    if err != nil {
        Reply(NewResponderRecordCreate(ctx, resource_str, nil, StatusFailed, errors.New(fmt.Sprintf("ParseJSON threw error: %s", err))));
    }

    if(ider == nil) {
        Reply(NewResponderRecordCreate(ctx, resource_str, nil, StatusFailed, errors.New("No error was thrown but ParseJSON did not return a valid object")));
    }
    if(rtype != nil && *rtype != resource_str) {
        Reply(NewResponderRecordCreate(ctx, resource_str, nil, StatusFailed, errors.New(fmt.Sprintf("This is resource \"%s\" but the new object includes type:\"%s\"", resource_str, rtype))));
    }
    if(id != nil && *id != "") {
        err = SetId(ider, *id);
        if(err != nil) {
            Reply(NewResponderRecordCreate(ctx, resource_str, nil, StatusFailed, errors.New(fmt.Sprintf("SetId failed:\"%s\"", err))));
        }
    }

    // first, we must check the permissions and verify that the
    // supplied linkages for each relationship is valid per the
    // rules of that relationship, eg, we don't want to let in
    // many linkages for a one to one relationship
    for _,linkage := range linkages.Linkages {
        fmt.Printf("Linkage: %#v\n", linkage);
        rel := a.RM.GetRelationship(resource_str, linkage.LinkName)
        fmt.Printf("REL: %s %s %#v\n", resource_str, linkage.LinkName, rel);
        if(rel == nil) {
            // user attempted to speify a relationship that does not exist
            panic("TODO: This");
        }
        rel.A.Authenticate("relationship.Create."+rel.SrcR+"."+rel.Name+"."+rel.DstR, "", r);
        err := rel.B.VerifyLinks(ider,linkage);
        if err != nil {
            Reply(NewResponderRecordCreate(ctx,resource_str, nil, StatusFailed, err));
        }
    }
    // trigger the pre-creates so the linkages have a chance to modify
    // the id object before it's inserted
    for _,linkage := range linkages.Linkages {
        rel := a.RM.GetRelationship(resource_str, linkage.LinkName)
        err := rel.B.PreCreate(ider,linkage);
        if err != nil {
            Reply(NewResponderRecordCreate(ctx,resource_str, nil, StatusFailed, err));
        }
    }

    createdStatus, err := resource.R.Create(ctx,resource_str,ider,id);
    if(err == nil && ider != nil && createdStatus & StatusCreated != 0) {
        for _,linkage := range linkages.Linkages {
            rel := a.RM.GetRelationship(resource_str, linkage.LinkName)
            err = rel.B.PostCreate(ider,linkage);
            if err != nil {
                Reply(NewResponderRecordCreate(ctx,resource_str, nil, StatusFailed, err));
            }
        }
    }
    Reply(NewResponderRecordCreate(ctx,resource_str, ider, createdStatus, err));
}

func(rr *RequestResolver) HandlerUpdate(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
