package jsonapi;

import ("net/http";);

type Worker interface {
    Work(a *API, r *http.Request);
    Cleanup(a *API, r *http.Request);
}

