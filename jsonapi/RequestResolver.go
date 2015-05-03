package jsonapi;

import ("net/http";"fmt";"strings";"github.com/julienschmidt/httprouter");

type RequestResolver struct{}

func NewRequestResolver() *RequestResolver {
    return &RequestResolver{};
}

func(rr *RequestResolver) HandlerFindOne(a *API, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    output := NewOutput(r);
    ids := strings.Split(ps.ByName("id"),",");
    res := []Ider{};
    var rmr *ResourceManagerResource;
    if len(ids) > 1 {
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
        fmt.Printf("Resource: %s\n", rmr.Name);
    }
    output.Data = NewOutputDataResources(false, data);
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
