package jsonapi;

//type RelationshipRequirement int;

// RelationshipBehavior is a "base interface"
// children: IdRelationshipBehavior or a HasIdRelationshipBehavior
type Relationship interface {
    IsSingle() bool
    //VerifyLinks(s Session, ider Ider, linkages *OutputLinkage) error
    //PreSave(s Session, ider Ider, linkages *OutputLinkage) error
    //PostSave(s Session, ider Ider, linkages *OutputLinkage) error
}

type RelationshipLinkId interface{
    Relationship
    //LinkId(s Session, srcR, dstR *ResourceManagerResource, src Ider) (ids []string)
}

type RelationshipLinkIder interface{
    Relationship
    //LinkIder(s Session, srcR, dstR *ResourceManagerResource,src Ider) (dst []Ider)
}

func VerifyRelationship(lb Relationship) bool {
    switch lb.(type) {
        case RelationshipLinkId:
            return true;
        case RelationshipLinkIder:
            return true;
        default:
            return false;
    }
}

