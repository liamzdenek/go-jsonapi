package jsonapi;

import (
    "net/http"
)

type Authenticator interface {
    Authenticate(permission, id string, r *http.Request)
}
