package jsonapi;

type Session interface {
    Begin() error
    Success() error
    Failure() error
}
