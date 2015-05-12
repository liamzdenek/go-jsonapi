package jsonapi;

import("fmt";"errors";"net/http")

type ResponderErrorRelationshipDoesNotExist struct {
    RelationshipName string
};

func NewResponderErrorRelationshipDoesNotExist(relname string) *ResponderErrorRelationshipDoesNotExist {
    return &ResponderErrorRelationshipDoesNotExist{
        RelationshipName: relname,
    }
}

func(e *ResponderErrorRelationshipDoesNotExist) Respond(a *API, w http.ResponseWriter, r *http.Request) error {
    re := NewResponderError(errors.New(
        fmt.Sprintf("The provided relationship \"%s\" does not exist.", e.RelationshipName),
    ));
    return re.Respond(a,w,r);
}
