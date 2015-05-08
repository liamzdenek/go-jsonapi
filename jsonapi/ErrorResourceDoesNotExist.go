package jsonapi;

import("fmt";"errors";"net/http")

type ErrorResourceDoesNotExist struct {
    ResourceName string
};

func NewErrorResourceDoesNotExist(relname string) *ErrorResourceDoesNotExist {
    return &ErrorResourceDoesNotExist{
        ResourceName: relname,
    }
}
func(e *ErrorResourceDoesNotExist) Respond(a *API, w http.ResponseWriter, r *http.Request) error {
    re := NewResponderError(errors.New(
        fmt.Sprintf("The provided resource \"%s\" does not exist.", e.ResourceName),
    ));
    return re.Respond(a,w,r);
}

