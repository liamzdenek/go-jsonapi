package jsonapi;

import( "net/http"; );

type Resource interface {
    FindOne(id string, r *http.Request) (Ider, error)
    FindMany(ids []string, r *http.Request) ([]Ider, error)
}
