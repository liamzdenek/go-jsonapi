package jsonapi;

import ("github.com/gorilla/pat";"net/http");

type API struct{
    Router *pat.Router
    BaseURIPath string
    RM *ResourceManager
}

func NewAPI() *API {
    api := &API{
        Router: pat.New(),
        RM: NewResourceManager(),
    };
    api.InitRouter();
    return api;
}

func (a *API) InitRouter() {

}

func (a *API) MountRelationship(name, srcR, dstR string, behavior RelationshipBehavior, auth Authenticator) {
    a.RM.MountRelationship(name,srcR,dstR,behavior,auth);
}

func (a *API) MountResource(name string, r Resource, auth Authenticator) {
    a.RM.MountResource(name,r,auth);
}

func(a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    //defer a.CatchResponses(w,r);
    a.Router.ServeHTTP(w, r);
}
