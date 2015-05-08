package jsonapi;

import("fmt";"errors";"net/http")

type ErrorOperationNotSupported struct {
    OperationDescription string
};

func NewErrorOperationNotSupported(desc string) *ErrorOperationNotSupported {
    return &ErrorOperationNotSupported{
        OperationDescription: desc,
    }
}

func(e *ErrorOperationNotSupported) Respond(a *API, w http.ResponseWriter, r *http.Request) error {
    re := NewResponderError(errors.New(
        fmt.Sprintf("The provided resource \"%s\" does not exist.", e.OperationDescription),
    ));
    return re.Respond(a,w,r);
}
