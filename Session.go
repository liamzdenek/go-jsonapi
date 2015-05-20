package jsonapi;

type Session interface {
    Begin(a *API) error
    Success(a *API) error
    Failure(a *API) error
}
