package jsonapi;

import ("net/http";);

type ResourceManagerRelationship struct {
    RM *ResourceManager
    SrcR string
    DstR string
    Name string
    B RelationshipBehavior
    A Authenticator
}

func(rmr *ResourceManagerRelationship) ResolveId(lb IdRelationshipBehavior, src Ider, generateIncluded bool) (*OutputLinkage, []IderTyper) {
    resource := rmr.RM.GetResource(rmr.DstR);
    res := &OutputLinkage{}
    included := []IderTyper{};
    ids := lb.Link(src);
    for _, id := range ids {
        res.Links = append(res.Links, OutputLink{
            Type: rmr.DstR,
            Id: id,
        });
    }
    if(generateIncluded) {
        linkdata, err := resource.R.FindMany(ids);
        Check(err);
        for _, link := range linkdata {
            included = append(included, NewIderTyperWrapper(link,rmr.DstR));
        }
    }
    return res, included;
}

func(rmr *ResourceManagerRelationship) ResolveIder(lb IderRelationshipBehavior, src Ider, generateIncluded bool) (*OutputLinkage, []IderTyper) {
    res := &OutputLinkage{}
    included := []IderTyper{};
    links := lb.Link(src);
    for _, link := range links {
        res.Links = append(res.Links, OutputLink{
            Type: rmr.DstR,
            Id: link.Id(),
        });
        //fixedlink,_ := a.AddLinkages(link, mr.DstR, r, false);
        //included = append(included, fixedlink);
    }
    return res, included;
}

func(rmr *ResourceManagerRelationship) Resolve(src Ider, r *http.Request, generateIncluded bool) (*OutputLinkage, []IderTyper) {
    // TODO: make this authentication request captured here
    rmr.A.Authenticate("relationship.FindAll."+rmr.SrcR+"."+rmr.Name+"."+rmr.DstR, src.Id(), r);
    // if we want included and it satisfies IderRelationshipBehavior, we 
    // should always prefer that over IdRelationshipBehavior
    if(generateIncluded) {
        if lb, found := rmr.B.(IderRelationshipBehavior); found {
            return rmr.ResolveIder(lb, src, generateIncluded);
        }
    }
    switch lb := rmr.B.(type) {
        case IdRelationshipBehavior:
            return rmr.ResolveId(lb, src, generateIncluded);
        case IderRelationshipBehavior:
            return rmr.ResolveIder(lb, src, generateIncluded);
        default:
            panic("Attempted to resolve a linkage behavior that is neither an Id or Ider LinkageBehavior.. This should never happen");
    }
}
