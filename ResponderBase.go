package jsonapi;

import ("net/http";)

func init() {
    var t Responder = &ResponderBase{};
    _ = t;
}

type ResponderBase struct{
    Output *Output
    Status int
    Headers map[string][]string
    CB func(a *API, s Session, r *http.Request);
}

func NewResponderBase(status int, o *Output) *ResponderBase {
    return &ResponderBase{
        Output: o,
        Status: status,
        Headers: map[string][]string{
            "Content-Type": []string{"application/vnd.api+json"},
        },
    }
}

func NewResponderBaseErrors(code int, es ...error) *ResponderBase {
    // TODO: replace this with code to make NewResponderErrors accept a list of OutputError to begin with
    oes := []OutputError{}
    for _,e := range es {
        oes = append(oes, OutputError{
            Title: e.Error(),
        });
    }
    o := NewOutput(nil);
    o.Errors = oes;
    return NewResponderBase(code, o);
}

func(rb *ResponderBase) PushHeader(k,v string) {
    if _, ok := rb.Headers[k]; !ok {
        rb.Headers[k] = []string{}
    }
    rb.Headers[k] = append(rb.Headers[k], v);
}

func(rb *ResponderBase) Respond(a *API, s Session, w http.ResponseWriter, r *http.Request) error {
    var err error;
    if !(rb.Status >= 200 && rb.Status < 300) || len(rb.Output.Errors) > 0 {
        err = s.Failure(a);
    } else {
        err = s.Success(a);
        if err != nil {
            err = s.Failure(a);
        }
    }
    if err != nil {
        panic(err); // TODO: properly encapsulate this into rrc.Output.Errors
    }
    
    rb.Output.Prepare();
    for k,vs := range rb.Headers {
        for _, v := range vs {
            w.Header().Add(k,v);
        }
    }
    w.WriteHeader(rb.Status);
    a.Send(rb.Output,w)
    return nil;
}
