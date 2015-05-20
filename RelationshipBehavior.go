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
    VerifyLinks(a *API, s Session, ider Ider, linkages *OutputLinkage) error
    PreCreate(a *API, s Session, ider Ider, linkages *OutputLinkage) error
    PostCreate(a *API, s Session, ider Ider, linkages *OutputLinkage) error
}

type IdRelationshipBehavior interface{
    RelationshipBehavior
    LinkId(a *API, s Session, srcR, dstR *ResourceManagerResource, src Ider) (ids []string)
}

type IderRelationshipBehavior interface{
    RelationshipBehavior
    LinkIder(a *API, s Session, srcR, dstR *ResourceManagerResource,src Ider) (dst []Ider)
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

