package jsonapi;

import ("net/http");

type Responder interface{
    Respond(*API, Session, http.ResponseWriter, *http.Request) error
}
