package jsonapi;

import ("net/http";);

type Task interface {
    Work(a *API, r *http.Request);
    ResponseWorker(has_paniced bool)
    Cleanup(a *API, r *http.Request);
}

