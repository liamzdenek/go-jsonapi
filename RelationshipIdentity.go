package jsonapi

type RelationshipIdentity struct {
	IsPrimary bool
}

func (ri *RelationshipIdentity) IsSingle() bool   { return true }
func (ri *RelationshipIdentity) PostMount(a *API) {}
func (ri *RelationshipIdentity) Link(r *Request, src, dst *ExecutableFuture, input FutureResponseKind) (output FutureRequestKind) {
	return &FutureRequestKindIdentity{
		Response: input,
		Future:   src.Future,
	}
}
func (ri *RelationshipIdentity) PushBackRelationships(r *Request, src, dst *ExecutableFuture, input, output FutureResponseKind) {
}
