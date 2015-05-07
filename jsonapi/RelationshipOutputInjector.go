package jsonapi;

import ("net/http";)

type RelationshipOutputInjector struct {
    Request *http.Request
    Ider Ider
    ResourceManagerResource *ResourceManagerResource
    A *API
    Include *IncludeInstructions
    Limit []string
}

func NewRelationshipOutputInjector(a *API, rmr *ResourceManagerResource, ider Ider, request *http.Request, include *IncludeInstructions) *RelationshipOutputInjector {
    return &RelationshipOutputInjector{
        A: a,
        ResourceManagerResource: rmr,
        Ider: ider,
        Request: request,
        Include: include,
    };
}

func(loi RelationshipOutputInjector) Link(included *[]Record) (*OutputLinkageSet) {
    rmr := loi.ResourceManagerResource;
    res := &OutputLinkageSet{
        RelatedBase: loi.A.GetBaseURL(loi.Request)+rmr.Name+"/"+loi.Ider.Id(),
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
