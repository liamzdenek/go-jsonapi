package jsonapi;

type SessionFactory interface {
    NewSession() Session
}
