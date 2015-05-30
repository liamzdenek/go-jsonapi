package jsonapi;

import ("net/http";);

type Task interface {
    Work(a *API, s Session, ctx *TaskContext, r *http.Request);
    ResponseWorker(has_paniced bool)
    Cleanup(a *API, r *http.Request);
}

