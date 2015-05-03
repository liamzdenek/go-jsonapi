package jsonapi;

type ResourceManagerRelationship struct {
    RM *ResourceManager
    SrcR string
    DstR string
    B RelationshipBehavior
    A Authenticator
}

func(rmr *ResourceManagerRelationship) Resolve(src Ider) *OutputLinkage {
    //resource := rmr.RM.GetResource(rmr.DstR);
    // TODO: perm check
    //rmr.A.Authenticate(mr.SrcR+".linkto."+mr.DstR+".FindMany", "", r);
    res := &OutputLinkage{}
    switch lb := rmr.B.(type) {
        case IdRelationshipBehavior:
            ids := lb.Link(src);
            for _, id := range ids {
                res.Links = append(res.Links, OutputLink{
                    Type: rmr.DstR,
                    Id: id,
                });
            }/*
            linkdata, err := resource.R.FindMany(ids);
            Check(err);
            for _, link := range linkdata {
                fixedlink,_ := a.AddLinkages(link, mr.DstR, r, false);
                included = append(included, fixedlink);
            }*/
        case IderRelationshipBehavior:
            links := lb.Link(src);
            for _, link := range links {
                res.Links = append(res.Links, OutputLink{
                    Type: rmr.DstR,
                    Id: link.Id(),
                });
                //fixedlink,_ := a.AddLinkages(link, mr.DstR, r, false);
                //included = append(included, fixedlink);
            }
        default:
            panic("Attempted to resolve a linkage behavior that is neither an Id or HasId LinkageBehavior.. This should never happen");
    }
    return res;
}
