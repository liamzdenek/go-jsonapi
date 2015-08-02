package jsonapi

type Responder interface {
	Respond(req *Request) error
}
