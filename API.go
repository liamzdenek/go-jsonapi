package jsonapi;

import (
    "github.com/julienschmidt/httprouter";
    "net/http"
);

/**
 * API is the primary user-facing structure within this framework. It
 * provides all of the functionality needed to intialize this framework,
 * as well as all of the glue to step down into more specific functionality
 */
type API struct {
    Resources map[string]APIMountedResource
    Relationships map[string]map[string]APIMountedRelationship
    Router *httprouter.Router
    Logger Logger
}

func NewAPI() *API {
    a := &API{
        Resources: map[string]APIMountedResource{},
        Relationships: map[string]map[string]APIMountedRelationship{},
        Router: httprouter.New(),
        Logger: NewLoggerDefault(nil),
    }
    a.InitRouter();
    return a;
}

/**
MountResource() will take a given Resource and make it available for requests sent to the given API. Any Resource that is accessible goes through this function
 */
func (a *API) MountResource(name string, resource Resource, authenticator Authenticator) {
    a.Resources[name] = APIMountedResource{
        Name: name,
        Resource: resource,
        Authenticator: authenticator,
    }
}

/**
MountRelationship() will take a given Relationship and make it available for requests sent to the given API. This also requires providing a source and destination Resource string. These resources must have already been mounted with MountResource() or this function will panic.
 */
func (a *API) MountRelationship(name, srcResourceName, dstResourceName string, relationship Relationship, authenticator Authenticator) {
    if _, exists := a.Resources[srcResourceName]; !exists {
        panic("Source resource "+srcResourceName+" for linkage does not exist");
    }
    if _, exists := a.Resources[dstResourceName]; !exists {
        panic("Destination resource "+dstResourceName+" for linkage does not exist");
    }
    if _, exists := a.Relationships[srcResourceName]; !exists {
        a.Relationships[srcResourceName] = make(map[string]APIMountedRelationship);
    }
    if(!VerifyRelationship(relationship)) {
        panic("Linkage provided cannot be used as an Id or Ider LinkageBehavior");
    }
    a.Relationships[srcResourceName][name] = APIMountedRelationship{
        SrcResourceName: srcResourceName,
        DstResourceName: dstResourceName,
        Name: name,
        Relationship: relationship,
        Authenticator: authenticator,
    };
}

/**
ServeHTTP() is to satisfy net/http.Handler -- all requests are simply forwarded through to httprouter
*/
func (a *API) ServeHTTP(w http.ResponseWriter,r *http.Request) {
    a.Router.ServeHTTP(w,r)
}

/**
InitRouter() prepares the internal httprouter object with all of the desired routes. This is called automatically. You should never have to call this unless you wish to muck around with the httprouter
*/
func (a *API) InitRouter() {
    a.Router.GET("/:resource/:id", a.Wrap(a.EntryFindRecordByResourceAndId));
}

/**
Wrap() reroutes a request to a standard httprouter.Handler (? double check) and converts it to the function signature that our entrypoint functions expect. It also initializes our panic handling and our thread pool handling.
*/
func(a *API) Wrap(child func(r *Request)) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
        req := NewRequest(a,r,w,params);
        defer req.Defer();
        child(req);
    }
}
