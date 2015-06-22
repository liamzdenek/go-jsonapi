package jsonapi;

type Promise interface {
    Success(r *Request)
    Failure(r *Request)
}
