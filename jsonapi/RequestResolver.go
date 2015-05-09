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
    output := NewOutput(r);
    ids := strings.Split(ps.ByName("id"),",");
    res := []Ider{};
    var rmr *ResourceManagerResource;
    resource_str := ps.ByName("resource");
    isSingle := len(ids) == 1;
    if !isSingle {
        res, rmr = rr.FindMany(a,r,resource_str,ids);
    } else {
        var tres Ider;
        tres, rmr = rr.FindOne(a,r,resource_str,ids[0]);
        res = []Ider{tres}
    }
    data := []*OutputDatum{};
    for _, ider := range res {
        roi := NewLinkerDefault(a, rmr, ider, r, ii);
        data = append(data, &OutputDatum{
            Datum: NewRecordWrapper(ider, rmr.Name, roi, true),
        });
    }
    fmt.Printf("Data: %#v\n", data);
    output.Data = NewOutputDataResources(isSingle, data);
    Reply(output);
}

func(rr *RequestResolver) FindOne(a *API, r *http.Request, resource_str, id_str string) (Ider, *ResourceManagerResource) {
    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(NewErrorResourceDoesNotExist(resource_str));
    }

    resource.A.Authenticate("resource.FindOne."+resource_str, id_str, r);

    data, err := resource.R.FindOne(id_str);
    Check(err);
    return data, resource;
}

func(rr *RequestResolver) FindMany(a *API, r *http.Request, resource_str string, ids []string) ([]Ider, *ResourceManagerResource) {
    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(NewErrorResourceDoesNotExist(resource_str));
    }

    id_str := strings.Join(ids, ",");

    resource.A.Authenticate("resource.FindMany."+resource_str, id_str, r);

    data, err := resource.R.FindMany(ids);
    Check(err);
    return data,resource;
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
    fmt.Printf("II: %#v\n",ii);
    output := NewOutput(r);
    ids := strings.Split(ps.ByName("id"),",");
    resource_str := ps.ByName("resource");
    var ider Ider;
    var rmr *ResourceManagerResource;
    if len(ids) > 1 {
        //res, rmr = rr.FindMany(a,r,ps,ids);
        // TODO: fix this maybe?
        panic(NewErrorOperationNotSupported("/:resource/:id/links does not support a list of links"));
    } else {
        ider, rmr = rr.FindOne(a,r,resource_str,ids[0]);
    }
    tii := NewIncludeInstructionsEmpty();
    tii.Push([]string{ps.ByName("linkname")});
    roi := NewLinkerDefault(a, rmr, ider, r, tii);
    wrapper := NewRecordWrapper(ider, rmr.Name, roi, true);

    include := &[]Record{};
    linkset := wrapper.Link(include);

    if(len(linkset.Linkages) == 0) {
        panic(NewErrorRelationshipDoesNotExist(ps.ByName("linkname")));
    }

    linkage := linkset.Linkages[0];

    if(len(linkage.Links) == 0) {
        // TODO: this
        output.Data.Data = nil;
        Reply(output);
    }
    if(len(linkage.Links) == 1) {
        lider, lrmr := rr.FindOne(a,r,linkage.Links[0].Type,linkage.Links[0].Id);
        
        // TODO: properly chain final argument here for includes
        lroi := NewLinkerDefault(a, lrmr, lider, r, ii);
        output.Data = NewOutputDataResources(true, []*OutputDatum{
            &OutputDatum{
                Datum: NewRecordWrapper(lider, lrmr.Name, lroi, true),
            },
        });
    } else {
        // TODO: it should
        panic(NewResponderError(errors.New("This request does not support one to many linkages")));
    }
    fmt.Printf("\nREPLYING\n\n");

    Reply(output);
}

/************************************************
 *
 * HandlerFindLinkByLinkNameAndResourceId is the entrypoint for
 * /:resource/:id/:linkname requests, primarily:
 * * /user/1/posts
 *
 * this handler does not support requests for multiple IDs:
 * * /user/1,2,3/posts
 *
 * requests with :linkname = "links" will 404
 */

func(rr *RequestResolver) HandlerFindLinkByNameAndResourceId(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    ii := NewIncludeInstructionsFromRequest(r);
    output := NewOutput(r);
    ids := strings.Split(ps.ByName("id"),",");
    resource_str := ps.ByName("resource");
    var ider Ider;
    var rmr *ResourceManagerResource;
    if len(ids) > 1 {
        panic(NewResponderError(errors.New("/:resource/:id/links/:linkname does not support a list of links")));
    } else {
        ider, rmr = rr.FindOne(a,r,resource_str,ids[0]);
    }
    roi := NewLinkerDefault(a, rmr, ider, r, ii);
    roi.Limit = []string{ps.ByName(":linkname")}
    wrapper := NewRecordWrapper(ider, rmr.Name, roi, true);

    include := &[]Record{};
    linkages := wrapper.Link(include);
    var link *OutputLinkage;
    if(linkages.Linkages != nil && len(linkages.Linkages) > 0) {
        link = linkages.Linkages[0];
    }

    output.Included = NewOutputIncluded(include);
    output.Data = NewOutputDataLinkage(true, link);
    Reply(output);
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
        panic(NewErrorResourceDoesNotExist(resource_str));
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
    resource_str := ps.ByName("resource");
    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(NewErrorResourceDoesNotExist(resource_str));
    }

    resource.A.Authenticate("resource.Create."+resource_str, "", r);

    body, err := ioutil.ReadAll(r.Body);
    if err != nil {
        panic(NewResponderError(errors.New(fmt.Sprintf("Body could not be parsed: %v\n", err))));
    }

    fmt.Printf("Request post body: %s\n", body);
    err = resource.R.Create(body);
    if err != nil {
        panic(NewResponderError(err));
    }
}
