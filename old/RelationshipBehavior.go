package jsonapi;

type RelationshipRequirement int;

const (
    Required RelationshipRequirement = iota;
    NotRequired
);

// RelationshipBehavior is a "base interface"
// children: IdRelationshipBehavior or a HasIdRelationshipBehavior
type RelationshipBehavior interface {
    IsSingle() bool
    VerifyLinks(s Session, ider Ider, linkages *OutputLinkage) error
    PreSave(s Session, ider Ider, linkages *OutputLinkage) error
    PostSave(s Session, ider Ider, linkages *OutputLinkage) error
}

type IdRelationshipBehavior interface{
    RelationshipBehavior
    LinkId(s Session, srcR, dstR *ResourceManagerResource, src Ider) (ids []string)
}

type IderRelationshipBehavior interface{
    RelationshipBehavior
    LinkIder(s Session, srcR, dstR *ResourceManagerResource,src Ider) (dst []Ider)
}

func VerifyRelationshipBehavior(lb RelationshipBehavior) bool {
    switch lb.(type) {
        case IdRelationshipBehavior:
            return true;
        case IderRelationshipBehavior:
            return true;
        default:
            return false;
    }
}

