package jsonapi;

import( "net/http"; );

type Resource interface {
    FindOne(id string, r *http.Request) (HasId, error)
    FindMany(ids []string, r *http.Request) ([]HasId, error)
}
