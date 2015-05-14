package jsonapi;

type Context interface {
    Begin() error
    Success() error
    Failure() error
}
