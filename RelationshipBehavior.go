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
    VerifyLinks(ider Ider, linkages *OutputLinkage) error
    PreCreate(ider Ider, linkages *OutputLinkage) error
    PostCreate(ider Ider, linkages *OutputLinkage) error
}

type IdRelationshipBehavior interface{
    LinkId(srcR, dstR *ResourceManagerResource, src Ider) (ids []string)
}

type IderRelationshipBehavior interface{
    LinkIder(srcR, dstR *ResourceManagerResource,src Ider) (dst []Ider)
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

