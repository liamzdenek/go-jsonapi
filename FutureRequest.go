package jsonapi

type FutureRequest struct {
	Request  *Request
	Response chan *FutureResponse
	Kind     FutureRequestKind
}

type FutureRequestedPanic struct{}

func NewFutureRequest(r *Request, kind FutureRequestKind) *FutureRequest {
	return &FutureRequest{
		Request:  r,
		Response: make(chan *FutureResponse),
		Kind:     kind,
	}
}

func (fr *FutureRequest) SendResponse(res *FutureResponse) {
	select {
	case fr.Response <- res:
	case <-fr.Request.Done.Wait():
		panic(&FutureRequestedPanic{})
	}
}

func (fr *FutureRequest) GetResponse() (res *FutureResponse) {
	select {
	case res = <-fr.Response:
	case <-fr.Request.Done.Wait():
		panic(&FutureRequestedPanic{})
	}
	return
}
