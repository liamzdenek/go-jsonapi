package jsonapi;

import ("github.com/gorilla/pat";"net/http";"fmt";"encoding/json");

type API struct{
    Router *pat.Router
    BaseURIPath string
    RM *ResourceManager
    RR *RequestResolver
}

func NewAPI() *API {
    api := &API{
        Router: pat.New(),
        RM: NewResourceManager(),
        RR: NewRequestResolver(),
    };
    api.InitRouter();
    return api;
}

// forwarding func to a.RM
func (a *API) MountRelationship(name, srcR, dstR string, behavior RelationshipBehavior, auth Authenticator) {
    a.RM.MountRelationship(name,srcR,dstR,behavior,auth);
}

// forwarding func to a.RM
func (a *API) MountResource(name string, r Resource, auth Authenticator) {
    a.RM.MountResource(name,r,auth);
}

// so the API can be mounted as a http handler
func(a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    defer a.CatchResponses(w,r);
    a.Router.ServeHTTP(w, r);
}

func(a *API) CatchResponses(w http.ResponseWriter, r *http.Request) {
    if raw := recover(); raw != nil {
        switch r := raw.(type) {
        case *Output:
            a.Send(r, w);
        case error:
            res := &Output{};
            res.Errors = []error{r};
            a.Send(res, w);
            panic(r);
        default:
            w.Write([]byte(fmt.Sprintf("Internal error handling request. Improper object sent to response method: %#v\n", r)));
        }
    }
}

func(a *API) Send(obj interface{}, w http.ResponseWriter) {
    str, err := json.Marshal(obj);
    Check(err);
    w.Write(str);
}



func (a *API) InitRouter() {
    a.Router.Get("/{resource}/{id}", a.Wrap(a.RR.HandlerFindOne));
}

func(a *API) Wrap(child func(a *API, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        child(a,w,r);
    }
}
