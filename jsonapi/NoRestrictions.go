package jsonapi;

import (
    "net/http"
)

type NoRestrictions struct {}

func NewNoRestrictions() *NoRestrictions {
    return &NoRestrictions{}
}

func(nr *NoRestrictions) Check(permission, id string, w http.ResponseWriter, r *http.Request) {

}
