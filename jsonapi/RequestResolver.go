package jsonapi;

import ("net/http";);

type RequestResolver struct{}

func NewRequestResolver() *RequestResolver {
    return &RequestResolver{};
}

func(rr *RequestResolver) HandlerFindOne(a *API, w http.ResponseWriter, r *http.Request) {

}
