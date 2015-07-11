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

type FutureResponseModifier interface{
    Modify(r *Request, src, dst *ExecutableFuture, k FutureResponseKind)
    EarlyModify(r *Request, src *ExecutableFuture)
}
type FutureResponseKind interface{}
type FutureResponseKindRecords struct{
    IsSingle bool
    Records []*Record
}

func(frr *FutureResponseKindRecords) EarlyModify(r *Request, src *ExecutableFuture) {
    for _, record := range frr.Records {
        if record.Type == "" {
            record.Type = src.Resource.Name;
        }
    }
    for _, rel := range r.API.GetRelationshipsByResource(src.Resource.Name) {
        for _, record := range frr.Records {
            newrel := &ORelationship{
                IsSingle: rel.Relationship.IsSingle(),
                RelationshipName: rel.Name,
                RelatedBase: r.GetBaseURL()+record.Type+"/"+record.Id,
            }
            record.PushRelationship(newrel);
        }
    }
}

func(frr *FutureResponseKindRecords) Modify(r *Request, src, dst *ExecutableFuture, rk FutureResponseKind) {
    switch k := rk.(type) {
    case *FutureResponseKindRecords:
        dstrel := &ORelationship{
            Data: GetResourceIdentifiers(frr.Records),
            RelationshipName: dst.Relationship.Name,
        };
        for _, record := range k.Records {
            record.PushRelationship(dstrel);
        }
        /*
        for _, record := range frr.Records {
            identifiers := GetResourceIdentifiers(k.Records);
            newrel := &ORelationship{
                IsSingle: k.IsSingle,
                Data: identifiers,
                RelationshipName: dst.Relationship.Name,
                RelatedBase: dst.Request.GetBaseURL(),
            };
            record.PushRelationship(newrel);
        }
        */
    default:
    }
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

type FutureRequestKindFindByFields struct{
    Fields []Field
}

