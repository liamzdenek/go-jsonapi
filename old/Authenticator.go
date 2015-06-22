package jsonapi;

type Authenticator interface {
    Authenticate(r *Request, permission, id string)
}
