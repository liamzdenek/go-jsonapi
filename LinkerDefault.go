package jsonapi;

import ("net/http";)

type LinkerDefault struct {
    Request *http.Request
    Ider Ider
    ResourceManagerResource *ResourceManagerResource
    A *API
    Include *IncludeInstructions
    Limit []string
}

func NewLinkerDefault(a *API, rmr *ResourceManagerResource, ider Ider, request *http.Request, include *IncludeInstructions) *LinkerDefault {
    return &LinkerDefault{
        A: a,
        ResourceManagerResource: rmr,
        Ider: ider,
        Request: request,
        Include: include,
    };
}

func(loi LinkerDefault) Link(included *[]Record) (*OutputLinkageSet) {
    rmr := loi.ResourceManagerResource;
    res := &OutputLinkageSet{
        RelatedBase: loi.A.GetBaseURL(loi.Request)+rmr.Name+"/"+GetId(loi.Ider),
    };
    for linkname,rel := range rmr.RM.GetRelationshipsByResource(rmr.Name){
        shouldFetch := loi.Include.ShouldFetch(linkname);
        link, new_included := rel.Resolve(loi.Ider, loi.Request, shouldFetch, loi.Include);
        link.LinkName = linkname;
        res.Linkages = append(res.Linkages, link);
        if(shouldFetch) {
            *included = append(*included, new_included...);
        }
    }
    return res;
}
