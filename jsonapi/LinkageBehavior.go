package jsonapi;

// LinkageBehavior is just syntax sugar.. what we're really
// looking for is either an IdLinkageBehavior or a HasIdLinkageBehavior
type LinkageBehavior interface {}

type IdLinkageBehavior interface{
    Link(src HasId) (ids []string)
}

type HasIdLinkageBehavior interface{
    Link(src HasId) (dst []HasId)
}

func VerifyLinkageBehavior(lb LinkageBehavior) bool {
    switch lb.(type) {
        case IdLinkageBehavior:
            return true;
        case HasIdLinkageBehavior:
            return true;
        default:
            return false;
    }
}

