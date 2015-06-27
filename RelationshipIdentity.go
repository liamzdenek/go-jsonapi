package jsonapi;

type RelationshipIdentity struct {
}

func(ri *RelationshipIdentity) IsSingle() bool { return true; }
func(ri *RelationshipIdentity) PostMount(a *API) {}
func(ri *RelationshipIdentity) Link(r *Request, src, dst *PreparedFuture, input FutureResponseKind) (output FutureRequestKind) {
    return &FutureRequestKindIdentity{
        Response: input,
        Future: src.Future,
    };
}
