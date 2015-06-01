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
    VerifyLinks(r *Request, rec *Record, rids []*OResourceIdentifier) error
    PreSave(r *Request, rec *Record, rids []*OResourceIdentifier) error
    PostSave(r *Request, rec *Record, rids []*OResourceIdentifier) error
}

type RelationshipLinkIds interface{
    Relationship
    LinkIds(r *Request, srcR, dstR *APIMountedResource, src *Record) (ids []string)
}

type RelationshipLinkRecords interface{
    Relationship
    LinkRecords(r *Request, srcR, dstR *APIMountedResource, src *Record) (dst []*Record)
}

func VerifyRelationship(lb Relationship) bool {
    switch lb.(type) {
        case RelationshipLinkIds:
            return true;
        case RelationshipLinkRecords:
            return true;
        default:
            return false;
    }
}
