package jsonapi;

import ("net/http";"fmt";"strings");

type RequestResolver struct{}

func NewRequestResolver() *RequestResolver {
    return &RequestResolver{};
}

func(rr *RequestResolver) HandlerFindOne(a *API, w http.ResponseWriter, r *http.Request) {
    output := NewOutput(r);
    res, rmr := rr.FindOne(a,r);
    include := strings.Split(r.URL.Query().Get("include"),",");
    roi := NewRelationshipOutputInjector(a, rmr, res, output, include);
    wrapped := NewIderLinkerTyperWrapper(res, rmr.Name, roi);
    fmt.Printf("Resource: %s\n", rmr.Name);
    output.Data = NewOutputDataResources(false, []*OutputDatum{
        &OutputDatum{
            Datum: wrapped,
        },
    });
    Reply(output);
}

func(rr *RequestResolver) FindOne(a *API, r *http.Request) (Ider, *ResourceManagerResource) {
    resource_str := r.URL.Query().Get(":resource");
    id_str := r.URL.Query().Get(":id");

    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(&ErrorResourceDoesNotExist{ResourceName:resource_str});
    }

    resource.A.Authenticate("resource.FindOne."+resource_str, id_str, r);

    data, err := resource.R.FindOne(id_str);
    Check(err);
    return data, resource;
}
