package jsonapi;

type Future interface{
    ShouldCombine(Future) bool
    Combine(Future) error
    Work(pf *PreparedFuture)
}

type FutureRequest struct {
    Request *Request
    Response chan *FutureResponse
    Kind FutureRequestKind
}

type FutureRequestedPanic struct{}

func NewFutureRequest(r *Request, kind FutureRequestKind) *FutureRequest {
    return &FutureRequest{
        Request: r,
        Response: make(chan *FutureResponse),
        Kind: kind,
    }
}

func(fr *FutureRequest) SendResponse(res *FutureResponse) {
    select {
    case fr.Response <- res:
    case <-fr.Request.Done.Wait():
        panic(&FutureRequestedPanic{});
    }
}


func(fr *FutureRequest) GetResponse() (res *FutureResponse){
    select {
    case res = <-fr.Response:
    case <-fr.Request.Done.Wait():
        panic(&FutureRequestedPanic{});
    }
    return
}

type FutureResponse struct {
    IsSuccess bool
    Success map[Future]FutureResponseKind
    Failure []OError
}

type FutureResponseKind interface{}
type FutureResponseKindRecords struct{
    IsSingle bool
    Records []*Record
}

type FutureRequestKind interface{}

type FutureRequestKindFailure struct {
    Response *FutureResponse
}
type FutureRequestKindIdentity struct {
    Response *FutureResponse
}
type FutureRequestKindFindByIds struct{
    Ids []string
}
type FutureRequestKindFindManyByField struct{}

