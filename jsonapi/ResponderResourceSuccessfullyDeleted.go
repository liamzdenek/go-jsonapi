package jsonapi;

import("net/http";);

type ResponderResourceSuccessfullyDeleted struct {
}

func NewResponderResourceSuccessfullyDeleted() *ResponderResourceSuccessfullyDeleted {
    return &ResponderResourceSuccessfullyDeleted{};
}

func(r *ResponderResourceSuccessfullyDeleted) Respond(a *API, w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(204); // No Content
}
