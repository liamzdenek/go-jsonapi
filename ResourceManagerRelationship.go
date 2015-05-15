package jsonapi;

import ("net/http";"fmt";"errors");

type ResourceManagerRelationship struct {
    RM *ResourceManager
    SrcR string
    DstR string
    Name string
    B RelationshipBehavior
    A Authenticator
    API *API
}

func(rmr *ResourceManagerRelationship) ResolveId(r *http.Request, lb IdRelationshipBehavior, src Ider, shouldFetch bool, include *IncludeInstructions) (*OutputLinkage, []Record) {
    resource := rmr.RM.GetResource(rmr.DstR);
    res := &OutputLinkage{}
    included := []Record{};
    dstRmr := rmr.RM.GetResource(rmr.DstR);
    ids := lb.LinkId(rmr.RM.GetResource(rmr.SrcR), dstRmr, src);
    for _, id := range ids {
        res.Links = append(res.Links, &OutputLink{
            Type: rmr.DstR,
            Id: id,
        });
    }
    if(shouldFetch) {
        linkdata, err := resource.R.FindMany(ids);
        Check(err);
        shouldInclude := include.ShouldInclude(rmr.Name);
        for _, link := range linkdata {
            fmt.Printf("Passing thru child include: %#v\n\n\n", include);
            roi := NewLinkerDefault(rmr.API, dstRmr, link, r, include.GetChild(rmr.Name));
            included = append(included, NewRecordWrapper(link,rmr.DstR,roi, shouldInclude));
        }
    }
    return res, included;
}

func(rmr *ResourceManagerRelationship) ResolveIder(r *http.Request, lb IderRelationshipBehavior, src Ider, shouldFetch bool, include *IncludeInstructions) (*OutputLinkage, []Record) {
    res := &OutputLinkage{}
    included := []Record{};
    dstRmr := rmr.RM.GetResource(rmr.DstR);
    links := lb.LinkIder(rmr.RM.GetResource(rmr.SrcR), dstRmr, src);
    shouldInclude := include.ShouldInclude(rmr.Name);
    for _, link := range links {
        res.Links = append(res.Links, &OutputLink{
            Type: rmr.DstR,
            Id: GetId(link),
        });
        fmt.Printf("\nShouldFetch %v ShouldInclude %v -- %s %#v\n\n", shouldFetch, shouldInclude, rmr.Name, include);
        if(shouldFetch) {
            roi := NewLinkerDefault(rmr.API, dstRmr, link, r, include.GetChild(rmr.Name));
            included = append(included, NewRecordWrapper(link,rmr.DstR, roi, shouldInclude));
        }
    }
    return res, included;
}

func(rmr *ResourceManagerRelationship) Resolve(src Ider, r *http.Request, shouldFetch bool, include *IncludeInstructions) (*OutputLinkage, []Record) {
    // TODO: make this authentication request captured here... a failure at a relationship should merely exclude that relationship
    rmr.A.Authenticate("relationship.FindAll."+rmr.SrcR+"."+rmr.Name+"."+rmr.DstR, GetId(src), r);
    // if we want included and it satisfies IderRelationshipBehavior, we 
    // should always prefer that over IdRelationshipBehavior
    if(shouldFetch) {
        if lb, found := rmr.B.(IderRelationshipBehavior); found {
            return rmr.ResolveIder(r, lb, src, shouldFetch, include);
        }
    }
    switch lb := rmr.B.(type) {
        case IdRelationshipBehavior:
            return rmr.ResolveId(r, lb, src, shouldFetch, include);
        case IderRelationshipBehavior:
            return rmr.ResolveIder(r, lb, src, shouldFetch, include);
        default:
            panic(NewResponderError(errors.New("Attempted to resolve a linkage behavior that is neither an Id or Ider LinkageBehavior.. This should never happen")));
    }
}
