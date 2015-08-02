package jsonapi

type Task interface {
	Work(r *Request)
	ResponseWorker(has_paniced bool)
	Cleanup(r *Request)
}
