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

func(e *ErrorRelationshipDoesNotExist) Respond(a *API, w http.ResponseWriter, r *http.Request) {
    re := NewResponderError(errors.New(
        fmt.Sprintf("The provided relationship \"%s\" does not exist.", e.RelationshipName),
    ));
    re.Respond(a,w,r);
}
