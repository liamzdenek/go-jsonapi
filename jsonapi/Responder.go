package jsonapi;

import ("net/http");

type Responder interface{
    Respond(*API, http.ResponseWriter, *http.Request) error
}
