package jsonapi;

import("fmt";"errors";"net/http")

type ResponderErrorOperationNotSupported struct {
    OperationDescription string
};

func NewResponderErrorOperationNotSupported(desc string) *ResponderErrorOperationNotSupported {
    return &ResponderErrorOperationNotSupported{
        OperationDescription: desc,
    }
}

func(e *ResponderErrorOperationNotSupported) Respond(a *API, w http.ResponseWriter, r *http.Request) error {
    re := NewResponderError(errors.New(
        fmt.Sprintf("The provided resource \"%s\" does not exist.", e.OperationDescription),
    ));
    return re.Respond(a,w,r);
}
