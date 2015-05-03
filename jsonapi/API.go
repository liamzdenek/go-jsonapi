package jsonapi;

import ("github.com/gorilla/pat";"net/http");

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
    //defer a.CatchResponses(w,r);
    a.Router.ServeHTTP(w, r);
}

func (a *API) InitRouter() {
    a.Router.Get("/{resource}/{id}", a.Wrap(a.RR.HandlerFindOne));
}

func(a *API) Wrap(child func(a *API, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        child(a,w,r);
    }
}
