package jsonapi;

type ContextFactory interface {
    NewContext() Context
}
