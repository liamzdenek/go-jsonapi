package jsonapi;

func init() {
    var t FutureResponseCanPushBackRelationships = &FutureResponseKindRecords{};
    t = &FutureResponseKindByFields{};
    _ = t;
}

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

type FutureResponseCanPushBackRelationships interface{
    PushBackRelationships(r *Request, src, dst *ExecutableFuture, k FutureResponseKind)
}
type FutureResponseKind interface{
}
type FutureResponseKindWithRecords interface{
    GetIsSingle() bool
    GetRecords() []*Record
}
type FutureResponseKindRecords struct{
    IsSingle bool
    Records []*Record
}

func(frkr *FutureResponseKindRecords) GetIsSingle() bool {return frkr.IsSingle;}
func(frkr *FutureResponseKindRecords) GetRecords() []*Record {return frkr.Records;}

func(frr *FutureResponseKindRecords) PushBackRelationships(r *Request, src, dst *ExecutableFuture, rk FutureResponseKind) {
    r.API.Logger.Debugf("GOT FUTURERESPONSEKIND: %#v\n", rk);;
    switch k := rk.(type) {
    case *FutureResponseKindRecords:
        dstrel := &ORelationship{
            Data: GetResourceIdentifiers(frr.Records),
            RelationshipName: dst.Relationship.Name,
        };
        for _, record := range k.Records {
            record.PushRelationship(dstrel);
        }
    case *FutureResponseKindByFields:
        for field, records := range k.Records {
            identifiers := GetResourceIdentifiers(records);
            for _, record := range frr.Records {
                if record.HasFieldValue(field) {
                    newrel := &ORelationship{
                        IsSingle: k.IsSingle,
                        Data: identifiers,
                        RelationshipName: dst.Relationship.Name,
                        RelatedBase: dst.Request.GetBaseURL(),
                    };
                    record.PushRelationship(newrel);
                }
            }
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
func(frr *FutureResponseKindByFields) PushBackRelationships(r *Request, src,dst *ExecutableFuture, rk FutureResponseKind) {
}

