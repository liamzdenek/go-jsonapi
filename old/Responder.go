package jsonapi;

import ("net/http");

type Responder interface{
    Respond(Session, http.ResponseWriter, *http.Request) error
}
