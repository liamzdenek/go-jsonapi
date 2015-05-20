package jsonapi;

import ("net/http";
"fmt";
"errors");

type ResourceManagerRelationship struct {
    RM *ResourceManager
    SrcR string
    DstR string
    Name string
    B RelationshipBehavior
    A Authenticator
    API *API
}

func(rmr *ResourceManagerRelationship) ResolveId(a *API, s Session, r *http.Request, lb IdRelationshipBehavior, src Ider, tctx *TaskContext, include *IncludeInstructions) (*OutputLinkage, []Record) {
    resource := rmr.RM.GetResource(rmr.DstR);
    res := &OutputLinkage{
        IsSingle: lb.IsSingle(),
    }
    included := []Record{};
    dstRmr := rmr.RM.GetResource(rmr.DstR);
    ids := lb.LinkId(a,s,rmr.RM.GetResource(rmr.SrcR), dstRmr, src);
    for _, id := range ids {
        res.Links = append(res.Links, &OutputLink{
            Type: rmr.DstR,
            Id: id,
        });
    }
    fmt.Printf("SHOULD FETCH LINK: %s %b\n\n", rmr.Name, include.ShouldFetch(rmr.Name));
    if(include.ShouldFetch(rmr.Name)) {
        linkdata, err := resource.R.FindMany(a, s, nil, ids);
        Check(err);
        for _, link := range linkdata {
            fmt.Printf("Passing thru child include: %#v\n\n\n", include);
            included = append(included, NewRecordWrapper(link,rmr.DstR,tctx, rmr.Name, include));
        }
    }
    return res, included;
}

func(rmr *ResourceManagerRelationship) ResolveIder(a *API, s Session, r *http.Request, lb IderRelationshipBehavior, src Ider, tctx *TaskContext, include *IncludeInstructions) (*OutputLinkage, []Record) {
    res := &OutputLinkage{
        IsSingle: lb.IsSingle(),
    }
    included := []Record{};
    dstRmr := rmr.RM.GetResource(rmr.DstR);
    links := lb.LinkIder(a,s,rmr.RM.GetResource(rmr.SrcR), dstRmr, src);
    for _, link := range links {
        res.Links = append(res.Links, &OutputLink{
            Type: rmr.DstR,
            Id: GetId(link),
        });
        included = append(included,
            NewRecordWrapper(link,rmr.DstR, tctx, rmr.Name, include),
        );
    }
    return res, included;
}

func(rmr *ResourceManagerRelationship) Resolve(a *API, s Session, src Ider, r *http.Request, shouldFetch bool, tctx *TaskContext, include *IncludeInstructions) (*OutputLinkage, []Record) {
    // TODO: make this authentication request captured here... a failure at a relationship should merely exclude that relationship
    rmr.A.Authenticate("relationship.FindAll."+rmr.SrcR+"."+rmr.Name+"."+rmr.DstR, GetId(src), r);
    // if we want included and it satisfies IderRelationshipBehavior, we 
    // should always prefer that over IdRelationshipBehavior
    if(shouldFetch) {
        if lb, found := rmr.B.(IderRelationshipBehavior); found {
            return rmr.ResolveIder(a, s, r, lb, src, tctx, include);
        }
    }
    switch lb := rmr.B.(type) {
        case IdRelationshipBehavior:
            return rmr.ResolveId(a, s, r, lb, src, tctx, include);
        case IderRelationshipBehavior:
            return rmr.ResolveIder(a, s, r, lb, src, tctx, include);
        default:
            panic(NewResponderError(errors.New("Attempted to resolve a linkage behavior that is neither an Id or Ider LinkageBehavior.. This should never happen")));
    }
}
