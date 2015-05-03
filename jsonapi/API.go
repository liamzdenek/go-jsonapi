package jsonapi;

import ("github.com/julienschmidt/httprouter";"net/http";"fmt";"encoding/json");

type API struct{
    Router *httprouter.Router
    BaseURIPath string
    RM *ResourceManager
    RR *RequestResolver
}

func NewAPI() *API {
    api := &API{
        Router: httprouter.New(),
        RM: NewResourceManager(),
        RR: NewRequestResolver(),
        BaseURIPath: "/",
    };
    api.InitRouter();
    return api;
}

func(a *API) GetBaseURL(r *http.Request) string {
    fmt.Printf("URL: %#v\n", r.URL);
    if r.URL.Scheme == "" {
        r.URL.Scheme = "http";
    }
    return r.URL.Scheme+"://"+r.Host+a.BaseURIPath;
}

// forwarding func to a.RM
func (a *API) MountRelationship(name, srcR, dstR string, behavior RelationshipBehavior, auth Authenticator) {
    a.RM.MountRelationship(name,srcR,dstR,behavior,auth);
}

// forwarding func to a.RM
func (a *API) MountResource(name string, r Resource, auth Authenticator) {
    a.RM.MountResource(name,r,auth);
}

// defines all the endpoints
func (a *API) InitRouter() {
    a.Router.GET("/:resource/:id", a.Wrap(a.RR.HandlerFindOne));
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

func(a *API) Wrap(child func(a *API, w http.ResponseWriter, r *http.Request, params httprouter.Params)) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
        child(a,w,r,params);
    }
}
