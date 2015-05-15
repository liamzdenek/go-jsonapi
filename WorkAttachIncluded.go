package jsonapi;

import("net/http");

type WorkAttachIncluded struct {
    Parent WorkerResultIders
}

func NewWorkAttachIncluded(parent WorkerResultIders) *WorkAttachIncluded {
    return &WorkAttachIncluded{
        Parent: parent,
    }
}

func (wai *WorkAttachIncluded) Work(a *API, r *http.Request) {
    
}

func (wai *WorkAttachIncluded) Cleanup(a *API, r *http.Request) {

}
