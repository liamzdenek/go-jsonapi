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
    VerifyLinks(r *Request, rec *Record, amr *APIMountedRelationship, rids []OResourceIdentifier) error
    PreSave(r *Request, rec *Record, amr *APIMountedRelationship, rids []OResourceIdentifier) error
    PostSave(r *Request, rec *Record, amr *APIMountedRelationship, rids []OResourceIdentifier) error
    Link(r *Request, srcR *APIMountedResource, rel *APIMountedRelationship, src Future) Future
    //PreDelete(r *Request, rec *Record, amr *APIMountedRelationship) error
    //PostDelete(r *Request, rec *Record, amr *APIMountedRelationship) error
}

