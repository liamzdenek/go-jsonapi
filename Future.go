package jsonapi;

type Future interface{
    ShouldCombine(Future) bool
    Combine(Future) error
    Work(requests chan *FutureRequest)
}

type FutureRequest struct {
    Request *Request
    Kind FutureRequestKind
    Response chan *FutureResponse
}

type FutureResponse struct {
    IsSuccess bool
    Success map[Future][]*Record
    Failure []*OError
}

type FutureRequestKind interface{}

type FutureRequestKindFindByIds struct{
    Ids []string
}
type FutureRequestKindFindManyByField struct{}

type FutureValue struct{
    Parent Future
    Field string
}
