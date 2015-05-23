package jsonapi;

import ("net/http";)

type ResponderErrors struct{
    Errors []error
}

func NewResponderError(e error) *ResponderErrors {
    return NewResponderErrors([]error{e});
}

func NewResponderErrors(e []error) *ResponderErrors {
    return &ResponderErrors{
        Errors: e,
    }
}

func(re *ResponderErrors) Respond(a *API, w http.ResponseWriter, r *http.Request) error {
    o := NewOutput(r,nil);
    o.Errors = re.Errors;
    o.Prepare();
    a.Send(o,w)
    return nil;
}
