package jsonapi;

type RelationshipRequirement int;

const (
    Required RelationshipRequirement = iota;
    NotRequired
);

// RelationshipBehavior is a "base interface"
// children: IdRelationshipBehavior or a HasIdRelationshipBehavior
type Relationship interface {
    IsSingle() bool
    PostMount(a *API)
    //VerifyLinks(r *Request, rec *Record, amr *APIMountedRelationship, rids []OResourceIdentifier) error
    //PreSave(r *Request, rec *Record, amr *APIMountedRelationship, rids []OResourceIdentifier) error
    //PostSave(r *Request, rec *Record, amr *APIMountedRelationship, rids []OResourceIdentifier) error
    //GetTargetFuture() Future
    Link(r *Request, src, dst *ExecutableFuture, input FutureResponseKind) (output FutureRequestKind)
}

