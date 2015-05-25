package jsonapi;

import ("net/http";)

type ResponderErrors struct{
    Errors []OutputError
}

func NewResponderError(e error) *ResponderErrors {
    return NewResponderErrors([]error{e});
}

func NewResponderErrors(es []error) *ResponderErrors {
    // TODO: replace this with code to make NewResponderErrors accept a list of OutputError to begin with
    oes := []OutputError{}
    for _,e := range es {
        oes = append(oes, OutputError{
            Title: e.Error(),
        });
    }
    return &ResponderErrors{
        Errors: oes,
    }
}

func(re *ResponderErrors) Respond(a *API, w http.ResponseWriter, r *http.Request) error {
    o := NewOutput(r,nil);
    o.Errors = re.Errors;
    o.Prepare();
    a.Send(o,w)
    return nil;
}
