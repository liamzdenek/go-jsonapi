package jsonapi

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/**
 * API is the primary user-facing structure within this framework. It
 * provides all of the functionality needed to intialize this framework,
 * as well as all of the glue to step down into more specific functionality
 */
type API struct {
	Resources              map[string]*APIMountedResource
	Relationships          map[string]map[string]*APIMountedRelationship
	DefaultResourceWrapper func(*APIMountedResource)
	Router                 *httprouter.Router
	BaseURI                string
	Logger                 Logger
}

func NewAPI(baseuri string) *API {
	a := &API{
		Resources:     map[string]*APIMountedResource{},
		Relationships: map[string]map[string]*APIMountedRelationship{},
		Router:        httprouter.New(),
		Logger:        NewLoggerDefault(nil),
		BaseURI:       baseuri,
		DefaultResourceWrapper: func(amr *APIMountedResource) {
			amr.Resource = amr.Resource
		},
	}
	a.InitRouter()
	return a
}

/**
MountResource() will take a given Resource and make it available for requests sent to the given API. Any Resource that is accessible goes through this function
*/
func (a *API) MountResource(name string, resource Resource, authenticators ...Authenticator) {
	amr := &APIMountedResource{
		Name:          name,
		Resource:      resource,
		Authenticator: NewAuthenticatorMany(authenticators...),
	}
	if a.DefaultResourceWrapper != nil {
		a.DefaultResourceWrapper(amr)
	}
	a.Resources[name] = amr
}

/**
MountRelationship() will take a given Relationship and make it available for requests sent to the given API. This also requires providing a source and destination Resource string. These resources must have already been mounted with MountResource() or this function will panic.
*/
func (a *API) MountRelationship(name, srcResourceName, dstResourceName string, relationship Relationship, authenticators ...Authenticator) {
	if _, exists := a.Resources[srcResourceName]; !exists {
		panic("Source resource " + srcResourceName + " for linkage does not exist")
	}
	if _, exists := a.Resources[dstResourceName]; !exists {
		panic("Destination resource " + dstResourceName + " for linkage does not exist")
	}
	if _, exists := a.Relationships[srcResourceName]; !exists {
		a.Relationships[srcResourceName] = make(map[string]*APIMountedRelationship)
	}
	amr := &APIMountedRelationship{
		SrcResourceName: srcResourceName,
		DstResourceName: dstResourceName,
		Name:            name,
		Relationship:    relationship,
		Authenticator:   NewAuthenticatorMany(authenticators...),
	}
	a.Relationships[srcResourceName][name] = amr
	amr.PostMount(a)
}

/**
GetResource() will return the resource for a given resource string. If the resource does not exist, this function returns a nil pointer.
*/
func (a *API) GetResource(name string) *APIMountedResource {
	r, ok := a.Resources[name]
	if !ok || r == nil {
		panic("The given resource " + name + " does not exist")
	}
	return r
}

/**
GetRelationship() will return a single relationship for a given resource string and relationship string. If the resource or relationship does not exist, this function returns a nil pointer.
*/
func (a *API) GetRelationship(srcR, linkName string) *APIMountedRelationship {
	if a.Resources[srcR] == nil {
		panic("The given resource " + srcR + " does not exist")
	}
	if a.Relationships[srcR] == nil || a.Relationships[srcR][linkName] == nil {
		panic("The given relationship " + srcR + "." + linkName + " does not exist")
	}
	return a.Relationships[srcR][linkName]
}

/**
GetRelationshipsByResource() will return a list of all of the relationships that the given resource string can link to.
*/
func (a *API) GetRelationshipsByResource(resource string) map[string]*APIMountedRelationship {
	return a.Relationships[resource]
}

/**
ServeHTTP() is to satisfy net/http.Handler -- all requests are simply forwarded through to httprouter
*/
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Router.ServeHTTP(w, r)
}

/**
InitRouter() prepares the internal httprouter object with all of the desired routes. This is called automatically. You should never have to call this unless you wish to muck around with the httprouter
*/
func (a *API) InitRouter() {
	// Resource and Relationship Read-only Operations
	a.Router.GET(a.BaseURI+":resource", a.Wrap(a.EntryFindDefault))
	a.Router.GET(a.BaseURI+":resource/:id", a.Wrap(a.EntryFindRecordByResourceAndId))
	a.Router.GET(a.BaseURI+":resource/:id/:linkname", a.Wrap(a.EntryFindRelationshipsByResourceId))
	a.Router.GET(a.BaseURI+":resource/:id/:linkname/:secondlinkname", a.Wrap(a.EntryFindRelationshipByNameAndResourceId))

	// Record Create
	a.Router.POST(a.BaseURI+":resource", a.Wrap(a.EntryCreate))
	a.Router.POST(a.BaseURI+":resource/:id", a.Wrap(a.EntryCreate))

	// Record Delete
	a.Router.DELETE(a.BaseURI+":resource/:id", a.Wrap(a.EntryDelete))
	//a.Router.DELETE(a.BaseURI+":resource", a.Wrap(a.EntryDelete));

	// Record Update
	a.Router.PATCH(a.BaseURI+":resource/:id", a.Wrap(a.EntryUpdate))
}

/**
Wrap() reroutes a request to a standard httprouter.Handler (? double check) and converts it to the function signature that our entrypoint functions expect. It also initializes our panic handling and our thread pool handling.
*/
func (a *API) Wrap(child func(r *Request)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		req := NewRequest(a, r, w, params)
		defer func() {
			if raw := recover(); raw != nil {
				req.HandlePanic(raw)
			}
			a.Logger.Debugf("CLLING REQ DONE WAIT\n")
			<-req.Done.Wait()
			a.Logger.Debugf("CALLING REQ DEFER\n")
			req.Defer()
		}()
		child(req)
	}
}
