package jsonapi;

import (
    "net/http"
)

type Authenticator interface {
    Authenticate(a *API, s Session, permission, id string, r *http.Request)
}
