package jsonapi;

import ("github.com/julienschmidt/httprouter";"net/http";"fmt";"encoding/json";"errors";"log";"os");

type API struct{
    Router *httprouter.Router
    BaseURIPath string
    RM *ResourceManager
    RR *RequestResolver
    SessionFactory
    Logger *log.Logger
}

func NewAPI(sf SessionFactory) *API {
    api := &API{
        Router: httprouter.New(),
        RM: NewResourceManager(),
        RR: NewRequestResolver(),
        BaseURIPath: "/",
        SessionFactory: sf,
        Logger: log.New(os.Stdout,"",log.LstdFlags | log.Lshortfile),
    };
    api.InitRouter();
    api.Logger.Printf("Initialized");
    return api;
}

func(a *API) GetBaseURL(r *http.Request) string {
    if r.URL.Scheme == "" {
        r.URL.Scheme = "http";
    }
    return r.URL.Scheme+"://"+r.Host+a.BaseURIPath;
}

// forwarding func to a.RM
func (a *API) MountRelationship(name, srcR, dstR string, behavior RelationshipBehavior, auth Authenticator) {
    a.RM.MountRelationship(name,srcR,dstR,behavior,auth,a);
}

// forwarding func to a.RM
func (a *API) MountResource(name string, r Resource, auth Authenticator) {
    a.RM.MountResource(name,r,auth);
}

// defines all the endpoints
func (a *API) InitRouter() {
    a.Router.GET("/:resource/:id/:linkname",
        a.WrapRedirector("linkname", "links",
            a.WrapPlain(http.NotFound), // if :linkname = "links"
            a.Wrap(a.RR.HandlerFindLinksByResourceId), // else
        ),
    );
    a.Router.GET("/:resource/:id/:linkname/:secondlinkname",
        a.WrapRedirector("linkname", "links",
            a.Wrap(a.RR.HandlerFindLinkByNameAndResourceId), // if :linkname = "links"
            a.WrapPlain(http.NotFound), // else
        ),
    );
    a.Router.DELETE("/:resource/:id", a.Wrap(a.RR.HandlerDelete));
    a.Router.PATCH("/:resource/:id", a.Wrap(a.RR.HandlerUpdate));
    a.Router.POST("/:resource", a.Wrap(a.RR.HandlerCreate));
    a.Router.GET("/:resource/:id", a.Wrap(a.RR.HandlerFindResourceById));
}

// so the API can be mounted as a http handler
func(a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    a.Router.ServeHTTP(w, r);
}

func(a *API) CatchResponses(w http.ResponseWriter, req *http.Request, raw interface{}) (was_handled bool, should_print_stack bool){
    a.Logger.Printf("Caught panic: %#v\n", raw);
    should_print_stack = true;
    switch r := raw.(type) {
    case Responder:
        should_print_stack = false;
        a.Logger.Printf("Respnding\n");
        r.Respond(a,w,req);
    case *Output:
        should_print_stack = false;
        a.Logger.Printf("Responder output\n");
        re := NewResponderOutput(r);
        re.Respond(a,w,req);
    case error:
        a.Logger.Printf("Panic(error) is deprecated as it is always ambiguous. You probably intend to use panic(NewResponderError()) instead\n");
        re := NewResponderError(r);
        re.Respond(a,w,req);
    case string:
        a.Logger.Printf("Panic(string) is deprecated as it is always ambiguous. You probably intend to use panic(NewResponderError()) instead\n");
        re := NewResponderError(errors.New(r));
        re.Respond(a,w,req);
    default:
        w.Write([]byte(fmt.Sprintf("Internal error handling request. Improper object sent to response method: %#v\n", r)));
        return false, true;
    }
    return true, should_print_stack;
}

func(a *API) Send(obj interface{}, w http.ResponseWriter) {
    str, err := json.Marshal(obj);
    Check(err);
    a.Logger.Printf("WRITING: %s\n", str);
    w.Write(str);
}

func(a *API) WrapRedirector(param_key, equal_to string, ifTrue httprouter.Handle, ifFalse httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
        if(p.ByName(param_key) == equal_to) {
            ifTrue(w,r,p);
        } else {
            ifFalse(w,r,p);
        }
    }
}

func(a *API) WrapPlain(child http.HandlerFunc) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
        child(w,r);
    }
}

func(a *API) Wrap(child func(a *API, w http.ResponseWriter, r *http.Request, params httprouter.Params)) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
        child(a,w,r,params);
    }
}

func (a *API) GetNewSession() Session {
    c := a.SessionFactory.NewSession();
    c.Begin();
    return c;
}
