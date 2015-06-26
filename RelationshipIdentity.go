package jsonapi;

type RelationshipIdentity struct {
}

func(ri *RelationshipIdentity) IsSingle() bool { return true; }
func(ri *RelationshipIdentity) PostMount(a *API) {}
func(ri *RelationshipIdentity) GetTargetFuture() Future { panic("RelationshipIdentity.GetTargetFuture should never be called"); }
func(ri *RelationshipIdentity) Link(r *Request, input *FutureResponse) (output FutureRequestKind) {
    return &FutureRequestKindIdentity{
        Response: input,
    };
}
