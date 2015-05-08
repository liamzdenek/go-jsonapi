package jsonapi;

import("fmt";"errors";"net/http")

type ErrorRelationshipDoesNotExist struct {
    RelationshipName string
};

func NewErrorRelationshipDoesNotExist(relname string) *ErrorRelationshipDoesNotExist {
    return &ErrorRelationshipDoesNotExist{
        RelationshipName: relname,
    }
}

func(e *ErrorRelationshipDoesNotExist) Respond(a *API, w http.ResponseWriter, r *http.Request) error {
    re := NewResponderError(errors.New(
        fmt.Sprintf("The provided relationship \"%s\" does not exist.", e.RelationshipName),
    ));
    return re.Respond(a,w,r);
}
