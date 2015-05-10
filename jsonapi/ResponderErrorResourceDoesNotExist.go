package jsonapi;

import("fmt";"errors";"net/http")

type ResponderErrorResourceDoesNotExist struct {
    ResourceName string
};

func NewResponderErrorResourceDoesNotExist(relname string) *ResponderErrorResourceDoesNotExist {
    return &ResponderErrorResourceDoesNotExist{
        ResourceName: relname,
    }
}
func(e *ResponderErrorResourceDoesNotExist) Respond(a *API, w http.ResponseWriter, r *http.Request) error {
    re := NewResponderError(errors.New(
        fmt.Sprintf("The provided resource \"%s\" does not exist.", e.ResourceName),
    ));
    return re.Respond(a,w,r);
}

