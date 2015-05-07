package jsonapi;

import ("net/http";);

type ResourceManagerRelationship struct {
    RM *ResourceManager
    SrcR string
    DstR string
    Name string
    B RelationshipBehavior
    A Authenticator
    API *API
}

func(rmr *ResourceManagerRelationship) ResolveId(r *http.Request, lb IdRelationshipBehavior, src Ider, generateIncluded bool, childInclude *IncludeInstructions) (*OutputLinkage, []Record) {
    resource := rmr.RM.GetResource(rmr.DstR);
    res := &OutputLinkage{}
    included := []Record{};
    dstRmr := rmr.RM.GetResource(rmr.DstR);
    ids := lb.LinkId(rmr.RM.GetResource(rmr.SrcR), dstRmr, src);
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
            roi := NewRelationshipOutputInjector(rmr.API, dstRmr, link, r, childInclude);
            included = append(included, NewRecordWrapper(link,rmr.DstR,roi));
        }
    }
    return res, included;
}

func(rmr *ResourceManagerRelationship) ResolveIder(r *http.Request, lb IderRelationshipBehavior, src Ider, generateIncluded bool, childInclude *IncludeInstructions) (*OutputLinkage, []Record) {
    res := &OutputLinkage{}
    included := []Record{};
    dstRmr := rmr.RM.GetResource(rmr.DstR);
    links := lb.LinkIder(rmr.RM.GetResource(rmr.SrcR), dstRmr, src);
    for _, link := range links {
        res.Links = append(res.Links, OutputLink{
            Type: rmr.DstR,
            Id: link.Id(),
        });
        //fixedlink,_ := a.AddLinkages(link, mr.DstR, r, false);
        if(generateIncluded) {
            roi := NewRelationshipOutputInjector(rmr.API, dstRmr, link, r, childInclude);
            included = append(included, NewRecordWrapper(link,rmr.DstR, roi));
        }
    }
    return res, included;
}

func(rmr *ResourceManagerRelationship) Resolve(src Ider, r *http.Request, generateIncluded bool, childInclude *IncludeInstructions) (*OutputLinkage, []Record) {
    // TODO: make this authentication request captured here... a failure at a relationship should merely exclude that relationship
    rmr.A.Authenticate("relationship.FindAll."+rmr.SrcR+"."+rmr.Name+"."+rmr.DstR, src.Id(), r);
    // if we want included and it satisfies IderRelationshipBehavior, we 
    // should always prefer that over IdRelationshipBehavior
    if(generateIncluded) {
        if lb, found := rmr.B.(IderRelationshipBehavior); found {
            return rmr.ResolveIder(r, lb, src, generateIncluded, childInclude);
        }
    }
    switch lb := rmr.B.(type) {
        case IdRelationshipBehavior:
            return rmr.ResolveId(r, lb, src, generateIncluded, childInclude);
        case IderRelationshipBehavior:
            return rmr.ResolveIder(r, lb, src, generateIncluded, childInclude);
        default:
            panic("Attempted to resolve a linkage behavior that is neither an Id or Ider LinkageBehavior.. This should never happen");
    }
}
