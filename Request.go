package jsonapi;

import (
    "net/http";
    "github.com/julienschmidt/httprouter";
    "errors";
    "fmt"
    "runtime"
    "encoding/json"
    "sync"
);

/**
Request is responsible for managing all of the common information between resources and relationships for the duration of a request. It contains references to often-needed components such as the raw net/http.Request, the API object, etc
*/
type Request struct {
    HttpRequest *http.Request
    HttpResponseWriter http.ResponseWriter
    API *API
    Params httprouter.Params
    IncludeInstructions *IncludeInstructions
    PromiseStorage *PromiseStorage
    hasCompleted bool
    Done *Done
    Responder chan Responder
}

type Done struct{
    sync.Mutex
    Chan chan bool 
    HasClosed bool
}

func NewDone() *Done {
    return &Done{
        Chan: make(chan bool),
    }
}

func(d *Done) Close() {
    d.Lock();
    defer d.Unlock();
    if !d.HasClosed {
        close(d.Chan);
    }
}

func(d *Done) Wait() chan bool {
    return d.Chan
}

/**
NewRequest() will return a populated instance of *Request. It will also initialize concurrency components.
*/
func NewRequest(a *API, httpreq *http.Request, httpres http.ResponseWriter, params httprouter.Params) * Request{
    req := &Request{
        API: a,
        HttpRequest: httpreq,
        HttpResponseWriter: httpres,
        Params: params,
        IncludeInstructions: NewIncludeInstructionsFromRequest(httpreq),
        PromiseStorage: NewPromiseStorage(),
        Responder: make(chan Responder),
        Done: NewDone(),
    }
    req.ResponderWorker();
    return req;
}

func(r *Request) ResponderWorker() {
    go func() {
        defer r.Done.Close();
        defer close(r.Responder);
        select {
        case <-r.Done.Wait():
        case re := <-r.Responder:
            re.Respond(r);
        }
    }()
}

func(r *Request) Respond(re Responder) {
    select {
    case r.Responder <- re:
    default:
    }
}

/**
Defer() should be called in a defer call at the same point that a Request is initialized. It is responsible for the safe handling of responses
*/
func(r *Request) Defer() {
    defer r.PromiseStorage.Defer();
    if(!r.hasCompleted) {
        panic("Did not call either Failure or Success before Defer");
    }
}

/**
Send() is responsible for converting a given *Output object to json, and sending it to the HttpResponseWriter that this Request is responsible for.
*/
func(r *Request) Send(obj interface{}) {
    str, err := json.Marshal(obj);
    Check(err);
    //r.API.Logger.Debugf("WRITING: %s TO %#v\n", str, r);
    r.HttpResponseWriter.Write(str);
}

func(r *Request) CatchPanic() {
    if raw := recover(); raw != nil {
        r.HandlePanic(raw);
    }
}

/**
HandlePanic() is responsible for interpreting the object that was paniced, and replying with the appropriate answer.
*/
func(r *Request) HandlePanic(raw interface{}){
    re,_ := r.InternalHandlePanic(raw);
    select {
    case r.Responder <- re:
    default:
    }
}

func(r *Request) InternalHandlePanic(raw interface{}) (re Responder, is_valid bool) {
    //r.API.Logger.Infof("Caught panic: %#v\n", raw);
    re, is_valid, should_print_stack := func() (Responder, bool, bool){
        switch raw_type := raw.(type) {
        case Responder:
            return raw_type, true, false;
        case *Output:
            rb := NewResponderBase(200,raw_type);
            return rb, true, false;
        case []OError:
            o := NewOutput();
            o.Errors = raw_type;
            re := NewResponderBase(500, o);
            return re, true, false;
        case error:
            r.API.Logger.Errorf("Panic(error) is deprecated as it is always ambiguous. You probably intend to use panic(NewResponderError()) instead\n");
            re := NewResponderBaseErrors(500, raw_type);
            return re, true, true;
        case string:
            r.API.Logger.Errorf("Panic(string) is deprecated as it is always ambiguous. You probably intend to use panic(NewResponderError()) instead\n");
            re := NewResponderBaseErrors(500,errors.New(raw_type));
            return re, true, true;
        default:
            r.API.Logger.Errorf("Panic(unknown type/%T) should be avoided as we cannot display a proper error to the user\n",raw_type);
            r.HttpResponseWriter.Write([]byte(fmt.Sprintf("Internal error handling request. Improper object sent to response method: %#v\n", r)));
            return nil, false, true;
        }
    }();
    if(should_print_stack) {
        const size = 64 << 10
        buf := make([]byte, size)
        buf = buf[:runtime.Stack(buf, false)]
        r.API.Logger.Infof("jsonapi: panic %#v\n%s", raw, buf);
    }
    return re,is_valid;

}
/**
Success() is responsible for calling the appropriate succcess handles. This function should never be called outside of a Responder 
*/
func(r *Request) Success() {
    //r.API.Logger.Infof("Calling Promise Success\n");
    r.finalizePromises(true);
}

/**
Failure() is responsible for calling the appropriate failure handles. This function should never be called outside of a Responder
*/
func(r *Request) Failure() {
    //r.API.Logger.Infof("Calling Promise Failure\n");
    r.finalizePromises(false);
}

func(r *Request) finalizePromises(success bool) {
    if(r.hasCompleted) {
        panic("Success or Failure can only be called once per request");
    }
    r.hasCompleted = true;
    for t,_ := range r.PromiseStorage.Promises {
        get := PromiseStorageLease{
            Type: t,
            ChanResponse: make(chan LeasedPromise),
        }
        r.PromiseStorage.ChanGet <- get;
        leased := <-get.ChanResponse;
        if(success) {
            leased.Promise.Success(r);
        } else {
            leased.Promise.Failure(r);
        }
        leased.Release();
    }
}

/**
GetBaseURL() will provide the URL + URI for any arbitrary request such that curling the output of this function is the root API endpoint for requests to this instance of this framework.
*/
func(r *Request) GetBaseURL() string {
    if r.HttpRequest.URL.Scheme == "" {
        r.HttpRequest.URL.Scheme = "http";
    }
    return r.HttpRequest.URL.Scheme+"://"+r.HttpRequest.Host+r.API.BaseURI;
}
