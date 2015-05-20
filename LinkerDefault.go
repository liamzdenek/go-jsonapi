package jsonapi;

import ("net/http";)

type LinkerDefault struct {
    Request *http.Request
    Ider Ider
    ResourceManagerResource *ResourceManagerResource
    A *API
    Include *IncludeInstructions
    Session 
    TaskContext *TaskContext
    Limit []string
}

func NewLinkerDefault(a *API, s Session, rmr *ResourceManagerResource, ider Ider, tctx *TaskContext, request *http.Request, include *IncludeInstructions) *LinkerDefault {
    return &LinkerDefault{
        A: a,
        ResourceManagerResource: rmr,
        Ider: ider,
        Request: request,
        Session: s,
        TaskContext: tctx,
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
        link, new_included := rel.Resolve(loi.A, loi.Session, loi.Ider, loi.Request, shouldFetch, loi.TaskContext,  loi.Include);
        link.LinkName = linkname;
        res.Linkages = append(res.Linkages, link);
        *included = append(*included, new_included...);
    }
    return res;
}
