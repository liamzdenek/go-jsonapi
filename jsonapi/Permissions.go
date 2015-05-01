package jsonapi;

import (
    "net/http"
)

type Permissions interface {
    Check(permission, id string, w http.ResponseWriter, r *http.Request)
}
