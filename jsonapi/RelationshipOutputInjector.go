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
    res := &OutputLinkageSet{};
    if relationships := rmr.RM.GetRelationshipsByResource(rmr.Name); len(relationships) > 0 {
        for linkname,rel := range relationships {
            link := rel.Resolve(loi.Ider);
            link.LinkName = linkname;
            res.Linkages = append(res.Linkages, link);
        }
    }
    return res;
}
