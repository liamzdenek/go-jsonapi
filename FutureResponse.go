package jsonapi

type FutureResponse struct {
	IsSuccess       bool
	Success         map[Future]FutureResponseKind
	Failure         []OError
	WaitForComplete chan bool
}
