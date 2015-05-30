package jsonapi;

import ("net/http");

/**
Request is responsible for managing all of the common information between resources and relationships for the duration of a request. It contains references to often-needed components such as the raw net/http.Request, the API object, etc
*/
type Request struct {
    HttpRequest *http.Request
    API *API
    //ThreadContext *ThreadContext
}
