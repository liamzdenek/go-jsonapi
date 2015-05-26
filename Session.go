package jsonapi;

import ("net/http")

type Session interface {
    SetData(d *SessionData)
    GetData() *SessionData
    Begin(a *API) error
    Success(a *API) error
    Failure(a *API) error
}

type SessionData struct {
    API *API
    Request *http.Request
}
