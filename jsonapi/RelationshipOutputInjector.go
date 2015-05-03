package jsonapi;

//import ("net/http");

type RelationshipOutputInjector struct {
    Output *Output
    Ider Ider
    ResourceManagerResource *ResourceManagerResource
    A *API
}

func NewRelationshipOutputInjector(a *API, rmr *ResourceManagerResource, ider Ider, output *Output) *RelationshipOutputInjector {
    return &RelationshipOutputInjector{
        A: a,
        ResourceManagerResource: rmr,
        Ider: ider,
        Output: output,
    };
}

func(loi RelationshipOutputInjector) Link() *OutputLinkageSet {
    rmr := loi.ResourceManagerResource;
    res := &OutputLinkageSet{
        RelatedBase: loi.A.GetBaseURL(loi.Output.Request)+rmr.Name+"/"+loi.Ider.Id(),
    };
    if relationships := rmr.RM.GetRelationshipsByResource(rmr.Name); len(relationships) > 0 {
        for linkname,rel := range relationships {
            link, included := rel.Resolve(loi.Ider, true);
            link.LinkName = linkname;
            res.Linkages = append(res.Linkages, link);
            loi.Output.Included.Push(included...);
        }
    }
    return res;
}
