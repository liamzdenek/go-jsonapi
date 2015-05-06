package jsonapi;

//import ("net/http");

type RelationshipOutputInjector struct {
    Output *Output
    Ider Ider
    ResourceManagerResource *ResourceManagerResource
    A *API
    Include *IncludeInstructions
    Limit []string
}

func NewRelationshipOutputInjector(a *API, rmr *ResourceManagerResource, ider Ider, output *Output, include *IncludeInstructions) *RelationshipOutputInjector {
    return &RelationshipOutputInjector{
        A: a,
        ResourceManagerResource: rmr,
        Ider: ider,
        Output: output,
        Include: include,
    };
}

func(loi *RelationshipOutputInjector) ShouldInclude(s string) bool {
    for _, include := range loi.Include.Instructions {
        if(include == s) {
            for _, limit := range loi.Limit {
                if(include == limit) {
                    return false;
                }
            }
            return true;
        }
    }
    return false;
}

func(loi RelationshipOutputInjector) Link(included *[]IderTyper) (*OutputLinkageSet) {
    rmr := loi.ResourceManagerResource;
    res := &OutputLinkageSet{
        RelatedBase: loi.A.GetBaseURL(loi.Output.Request)+rmr.Name+"/"+loi.Ider.Id(),
    };
    if relationships := rmr.RM.GetRelationshipsByResource(rmr.Name); len(relationships) > 0 {
        for linkname,rel := range relationships {
            shouldInclude := loi.ShouldInclude(linkname);
            link, new_included := rel.Resolve(loi.Ider, loi.Output.Request, shouldInclude);
            link.LinkName = linkname;
            res.Linkages = append(res.Linkages, link);
            if(shouldInclude) {
                *included = append(*included, new_included...);
            }
        }
    }
    return res;
}
