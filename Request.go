package jsonapi;

import (
    "net/http";
    "github.com/julienschmidt/httprouter";
    "errors";
    "fmt"
    "runtime"
    "encoding/json"
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
    TaskContext *TaskContext
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
    }
    req.TaskContext = NewTaskContext(req);
    return req;
}

/**
Defer() should be called in a defer call at the same point that a Request is initialized. It is responsible for the safe handling of responses
*/
func(r *Request) Defer() {
    defer r.TaskContext.Cleanup();
}

/**
Send() is responsible for converting a given *Output object to json, and sending it to the HttpResponseWriter that this Request is responsible for.
*/
func(r *Request) Send(obj interface{}) {
    str, err := json.Marshal(obj);
    Check(err);
    r.API.Logger.Debugf("WRITING: %s\n", str);
    r.HttpResponseWriter.Write(str);
}

/**
HandlePanic() is responsible for interpreting the object that was paniced, and replying with the appropriate answer.
*/
func(r *Request) HandlePanic(raw interface{}) (is_valid bool){
    r.API.Logger.Infof("Caught panic: %#v\n", raw);
    is_valid, should_print_stack := func() (bool, bool){
        switch raw_type := raw.(type) {
        case Responder:
            r.API.Logger.Infof("Responding\n");
            raw_type.Respond(r);
            return true, false;
        case *Output:
            r.API.Logger.Infof("Responder output\n");
            rb := NewResponderBase(200,raw_type);
            rb.Respond(r);
            return true, false;
        case error:
            r.API.Logger.Errorf("Panic(error) is deprecated as it is always ambiguous. You probably intend to use panic(NewResponderError()) instead\n");
            re := NewResponderBaseErrors(500, raw_type);
            re.Respond(r);
            return true, true;
        case string:
            r.API.Logger.Errorf("Panic(string) is deprecated as it is always ambiguous. You probably intend to use panic(NewResponderError()) instead\n");
            re := NewResponderBaseErrors(500,errors.New(raw_type));
            re.Respond(r);
            return true, true;
        default:
            r.HttpResponseWriter.Write([]byte(fmt.Sprintf("Internal error handling request. Improper object sent to response method: %#v\n", r)));
            return false, true;
        }
    }();
    if(should_print_stack) {
        const size = 64 << 10
        buf := make([]byte, size)
        buf = buf[:runtime.Stack(buf, false)]
        r.API.Logger.Infof("jsonapi: panic %v\n%s", raw, buf);
    }
    return is_valid;
}

/**
Success() is responsible for calling the appropriate succcess handles. This function should never be called outside of a Responder 
*/
func(r *Request) Success() {
    //TODO: hook it up
}

/**
Failure() is responsible for calling the appropriate failure handles. This function should never be called outside of a Responder
*/
func(r *Request) Failure() {
    // TODO: hook it up
}
