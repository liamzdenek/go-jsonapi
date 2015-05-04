package jsonapi;

import ("net/http";"strings";"github.com/julienschmidt/httprouter");

type RequestResolver struct{}

func NewRequestResolver() *RequestResolver {
    return &RequestResolver{};
}

/************************************************
 *
 * HandlerFindOne  is the entrypoint for /:resource/:id requests, primarily:
 * * /user/1
 * * /user/1,2,3
 */

func(rr *RequestResolver) HandlerFindOne(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    output := NewOutput(r);
    ids := strings.Split(ps.ByName("id"),",");
    res := []Ider{};
    var rmr *ResourceManagerResource;
    isSingle := len(ids) == 1;
    if !isSingle {
        res, rmr = rr.FindMany(a,r,ps,ids);
    } else {
        var tres Ider;
        tres, rmr = rr.FindOne(a,r,ps,ids[0]);
        res = []Ider{tres}
    }
    data := []*OutputDatum{};
    for _, ider := range res {
        include := strings.Split(r.URL.Query().Get("include"),",");
        roi := NewRelationshipOutputInjector(a, rmr, ider, output, include);
        data = append(data, &OutputDatum{
            Datum: NewIderLinkerTyperWrapper(ider, rmr.Name, roi),
        });
    }
    output.Data = NewOutputDataResources(isSingle, data);
    Reply(output);
}

func(rr *RequestResolver) FindOne(a *API, r *http.Request, ps httprouter.Params, id_str string) (Ider, *ResourceManagerResource) {
    resource_str := ps.ByName("resource");

    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(&ErrorResourceDoesNotExist{ResourceName:resource_str});
    }

    resource.A.Authenticate("resource.FindOne."+resource_str, id_str, r);

    data, err := resource.R.FindOne(id_str);
    Check(err);
    return data, resource;
}

func(rr *RequestResolver) FindMany(a *API, r *http.Request, ps httprouter.Params, ids []string) ([]Ider, *ResourceManagerResource) {
    resource_str := ps.ByName("resource");
    id_str := ps.ByName("id");

    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(&ErrorResourceDoesNotExist{ResourceName:resource_str});
    }

    resource.A.Authenticate("resource.FindMany."+resource_str, id_str, r);

    data, err := resource.R.FindMany(ids);
    Check(err);
    return data,resource;
}

/************************************************
 *
 * HandlerFindOneLinks  is the entrypoint for /:resource/:id requests, primarily:
 * * /user/1/links
 *
 * this handler does not support requests for multiple IDs (maybe it could?):
 * * /user/1,2,3/links
 */

func(rr *RequestResolver) HandlerFindOneLinks(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    output := NewOutput(r);
    ids := strings.Split(ps.ByName("id"),",");
    var ider Ider;
    var rmr *ResourceManagerResource;
    if len(ids) > 1 {
        //res, rmr = rr.FindMany(a,r,ps,ids);
        // TODO: fix this maybe?
        panic("/:resource/:id/links does not support a list of links");
    } else {
        ider, rmr = rr.FindOne(a,r,ps,ids[0]);
    }
    include := strings.Split(r.URL.Query().Get("include"),",");
    roi := NewRelationshipOutputInjector(a, rmr, ider, output, include);
    wrapper := NewIderLinkerTyperWrapper(ider, rmr.Name, roi);

    output.Data = NewOutputDataRelationship(wrapper.Link());
    Reply(output);
}

/************************************************
 *
 * HandlerFindOneSpecificLink  is the entrypoint for /:resource/:id requests, primarily:
 * * /user/1/links/posts
 *
 * this handler does not support requests for multiple IDs (maybe it could?):
 * * /user/1,2,3/links/posts
 */

func(rr *RequestResolver) HandlerFindOneSpecificLink(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    output := NewOutput(r);
    ids := strings.Split(ps.ByName("id"),",");
    var ider Ider;
    var rmr *ResourceManagerResource;
    if len(ids) > 1 {
        //res, rmr = rr.FindMany(a,r,ps,ids);
        // TODO: fix this maybe?
        panic("/:resource/:id/links/:linkname does not support a list of links");
    } else {
        ider, rmr = rr.FindOne(a,r,ps,ids[0]);
    }
    include := strings.Split(r.URL.Query().Get("include"),",");
    roi := NewRelationshipOutputInjector(a, rmr, ider, output, include);
    roi.Limit = []string{ps.ByName(":linkname")}
    wrapper := NewIderLinkerTyperWrapper(ider, rmr.Name, roi);

    linkages := wrapper.Link();
    var link *OutputLinkage;
    if(linkages.Linkages != nil) {
        link = linkages.Linkages[0];
    }

    output.Data = NewOutputDataLinkage(true, link);
    Reply(output);
}

