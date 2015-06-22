package jsonapi;

import("errors";)

type APIMountedRelationship struct {
    SrcResourceName string
    Name string
    Relationship
    Authenticator
}

func(amr *APIMountedRelationship) Resolve(r *Request, src *Record, shouldFetch bool, include *IncludeInstructions) (*ORelationship, []*Record) {
    amr.Authenticator.Authenticate(r,"relationship.FindAll."+amr.SrcResourceName+"."+amr.Name, src.Id);
    if lb, found := amr.Relationship.(RelationshipLinkRecords); shouldFetch && found {
        return amr.ResolveRecords(r, lb, src, include);
    }
    switch lb := amr.Relationship.(type) {
        case RelationshipLinkIds:
            return amr.ResolveIds(r, lb, src, include);
        case RelationshipLinkRecords:
            return amr.ResolveRecords(r, lb, src, include);
        default:
            panic(NewResponderBaseErrors(500, errors.New("Attempted to resolve a linkage behavior that is neither an Id or Ider LinkageBehavior.. This should never happen")));
    }
}

func(amr *APIMountedRelationship) ResolveIds(r *Request, lb RelationshipLinkIds, record *Record, include *IncludeInstructions) (*ORelationship, []*Record) {
    rel := &ORelationship{
        IsSingle: lb.IsSingle(),
    }
    included := []*Record{};
    srcResource := r.API.GetResource(amr.SrcResourceName);
    ids := lb.LinkIds(r,srcResource,amr,record);
    shouldFetch := include.ShouldFetch(amr.Name);
    for _, id := range ids {
        rel.Data = append(rel.Data, id);
        if(shouldFetch) {
            dstResource := r.API.GetResource(id.Type)
            one, err := dstResource.FindOne(r, RequestParams{}, id.Id);
            if err != nil {
                panic(err);
            }
            one.ShouldInclude = include.ShouldInclude(amr.Name);
            r.API.Logger.Debugf("SHOULD INCLUDE: %s %v\n", amr.Name, one.ShouldInclude);
            included = append(included, one);
        }
    }
    return rel,included;
}

func(amr *APIMountedRelationship) ResolveRecords(r *Request, lb RelationshipLinkRecords, record *Record, include *IncludeInstructions) (*ORelationship, []*Record) {
    rel := &ORelationship{
        IsSingle: lb.IsSingle(),
    };
    included := []*Record{};
    srcResource := r.API.GetResource(amr.SrcResourceName);
    records := lb.LinkRecords(r,srcResource,amr,record);
    for _, record := range records {
        record.PrepareRelationships(r, include.GetChild(amr.Name));
        rel.Data = append(rel.Data, record.GetResourceIdentifier());
        record.ShouldInclude = include.ShouldInclude(amr.Name);
        r.API.Logger.Debugf("SHOULD INCLUDE: %s %v\n", amr.Name, record.ShouldInclude);
        included = append(included, record);
    }
    return rel,included;
}
