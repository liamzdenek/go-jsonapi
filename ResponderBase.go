package jsonapi;

func init() {
    var t Responder = &ResponderBase{};
    _ = t;
}

type ResponderBase struct{
    Output *Output
    Status int
    Headers map[string][]string
    CB func(r *Request);
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
    o := NewOutput();
    o.Errors = oes;
    return NewResponderBase(code, o);
}

func(rb *ResponderBase) PushHeader(k,v string) {
    if _, ok := rb.Headers[k]; !ok {
        rb.Headers[k] = []string{}
    }
    rb.Headers[k] = append(rb.Headers[k], v);
}

func(rb *ResponderBase) Respond(r *Request) error {
    // TODO: improve this?
    if !(rb.Status >= 200 && rb.Status < 300) || (rb.Output != nil && rb.Output.Errors != nil && len(rb.Output.Errors) > 0) {
        r.Failure();
    } else {
        r.Success();
    }
    for k,vs := range rb.Headers {
        for _, v := range vs {
            r.HttpResponseWriter.Header().Add(k,v);
        }
    }
    r.HttpResponseWriter.WriteHeader(rb.Status);
    if rb.Output != nil {
        r.Send(rb.Output)
    }
    return nil;
}
