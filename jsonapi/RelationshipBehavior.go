package jsonapi;

// RelationshipBehavior is just syntax sugar.. what we're really
// looking for is either an IdRelationshipBehavior or a HasIdRelationshipBehavior
type RelationshipBehavior interface {}

type IdRelationshipBehavior interface{
    Link(src Ider) (ids []string)
}

type IderRelationshipBehavior interface{
    Link(src Ider) (dst []Ider)
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

