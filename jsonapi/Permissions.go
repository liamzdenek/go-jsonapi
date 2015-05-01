package jsonapi;

import (
    "net/http"
)

type Permissions interface {
    Check(permission, id string, r *http.Request)
}
