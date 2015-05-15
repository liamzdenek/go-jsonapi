package jsonapi;

import ("net/http";);

type Worker interface {
    Work(a *API, r *http.Request);
    ResponseWorker(has_paniced bool)
    Cleanup(a *API, r *http.Request);
}

