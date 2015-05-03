package jsonapi;

import ("net/http";"fmt");

type RequestResolver struct{}

func NewRequestResolver() *RequestResolver {
    return &RequestResolver{};
}

func(rr *RequestResolver) HandlerFindOne(a *API, w http.ResponseWriter, r *http.Request) {
    res, resource_str := rr.FindOne(a,r);
    wrapped := NewIderLinkerTyperWrapper(res, resource_str)
    fmt.Printf("Resource: %s\n", resource_str);
    Reply(&Output{
        Data: NewOutputDataResources(false, []*OutputDatum{
            &OutputDatum{
                Datum: wrapped,
            },
        }),
    });
}

func(rr *RequestResolver) FindOne(a *API, r *http.Request) (Ider, string) {
    resource_str := r.URL.Query().Get(":resource");
    id_str := r.URL.Query().Get(":id");
    
    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(&ErrorResourceDoesNotExist{ResourceName:resource_str});
    }

    resource.A.Authenticate(resource_str+".FindAll", id_str, r);

    data, err := resource.R.FindOne(id_str, r);
    Check(err);
    return data, resource_str;
}
