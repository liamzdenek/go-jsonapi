package jsonapi;

import ("net/http";"strings";"github.com/julienschmidt/httprouter";"fmt");

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
    fmt.Printf("II: %#v\n",ii);
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
        roi := NewRelationshipOutputInjector(a, rmr, ider, output, ii);
        data = append(data, &OutputDatum{
            Datum: NewIderLinkerTyperWrapper(ider, rmr.Name, roi),
        });
    }
    output.Data = NewOutputDataResources(isSingle, data);
    Reply(output);
}

func(rr *RequestResolver) FindOne(a *API, r *http.Request, resource_str, id_str string) (Ider, *ResourceManagerResource) {
    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(&ErrorResourceDoesNotExist{ResourceName:resource_str});
    }

    resource.A.Authenticate("resource.FindOne."+resource_str, id_str, r);

    data, err := resource.R.FindOne(id_str);
    Check(err);
    return data, resource;
}

func(rr *RequestResolver) FindMany(a *API, r *http.Request, resource_str string, ids []string) ([]Ider, *ResourceManagerResource) {
    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(&ErrorResourceDoesNotExist{ResourceName:resource_str});
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
    //ii := NewIncludeInstructionsFromRequest(r);
    output := NewOutput(r);
    ids := strings.Split(ps.ByName("id"),",");
    resource_str := ps.ByName("resource");
    var ider Ider;
    var rmr *ResourceManagerResource;
    if len(ids) > 1 {
        //res, rmr = rr.FindMany(a,r,ps,ids);
        // TODO: fix this maybe?
        panic("/:resource/:id/links does not support a list of links");
    } else {
        ider, rmr = rr.FindOne(a,r,resource_str,ids[0]);
    }
    tii := NewIncludeInstructionsEmpty();
    tii.Push([]string{ps.ByName("linkname")});
    roi := NewRelationshipOutputInjector(a, rmr, ider, output, tii);
    wrapper := NewIderLinkerTyperWrapper(ider, rmr.Name, roi);

    include := &[]IderTyper{};
    linkset := wrapper.Link(include);

    if(len(linkset.Linkages) == 0) {
        // TODO: spec compliance
        panic("This linkage does not exist");
    }

    linkage := linkset.Linkages[0];

    if(len(linkage.Links) == 0) {
        // TODO: this
        panic("this should return with primary data as null");
    }

    if(len(linkage.Links) == 1) {
        lider, lrmr := rr.FindOne(a,r,linkage.Links[0].Type,linkage.Links[0].Id);

        // TODO: properly chain final argument here for includes
        lroi := NewRelationshipOutputInjector(a, lrmr, lider, output, NewIncludeInstructionsEmpty());
        output.Data = NewOutputDataResources(true, []*OutputDatum{
            &OutputDatum{
                Datum: NewIderLinkerTyperWrapper(lider, lrmr.Name, lroi),
            },
        });
    } else {

    }

    fmt.Printf("Linkset: %#v\n", linkset.Linkages[0]);
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
        panic("/:resource/:id/links/:linkname does not support a list of links");
    } else {
        ider, rmr = rr.FindOne(a,r,resource_str,ids[0]);
    }
    roi := NewRelationshipOutputInjector(a, rmr, ider, output, ii);
    roi.Limit = []string{ps.ByName(":linkname")}
    wrapper := NewIderLinkerTyperWrapper(ider, rmr.Name, roi);

    include := &[]IderTyper{};
    linkages := wrapper.Link(include);
    var link *OutputLinkage;
    if(linkages.Linkages != nil) {
        link = linkages.Linkages[0];
    }

    output.Data = NewOutputDataLinkage(true, link);
    output.Included = NewOutputIncluded(*include);
    Reply(output);
}

