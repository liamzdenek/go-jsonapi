package jsonapi;

type Future interface{
    ShouldCombine(Future) bool
    Combine(Future) error
    Work(ef *ExecutableFuture)
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
    WaitForComplete chan bool
}

type FutureResponseKind interface{
}
type FutureResponseKindWithRecords interface{
    GetIsSingle() bool
    GetRecords() []*Record
}

type FutureRequestKind interface{}

type FutureRequestKindFailure struct {
    Response *FutureResponse
}
type FutureRequestKindIdentity struct {
    Response FutureResponseKind
    Future
}
type FutureRequestKindFindByIds struct{
    Ids []string
}

type Field struct {
    Field string
    Value string
}

type FutureRequestKindFindByAnyFields struct{
    Fields []Field
}

type FutureResponseKindByFields struct{
    IsSingle bool
    Records map[Field][]*Record
}
func(frkbf *FutureResponseKindByFields) GetIsSingle() bool {
    return false;
}
func(frkbf *FutureResponseKindByFields) GetRecords() []*Record {
    rec := []*Record{};
    for _, recs := range frkbf.Records {
        rec = append(rec, recs...);
    }
    return rec;
}
