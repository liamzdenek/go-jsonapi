package jsonapi;

import("errors";)

type APIMountedRelationship struct {
    SrcResourceName string
    DstResourceName string
    Name string
    Relationship
    Authenticator
}

func(amr *APIMountedRelationship) Resolve(r *Request, src *Record, shouldFetch bool, include *IncludeInstructions) (*ORelationship, []*Record) {
    amr.Authenticator.Authenticate(r,"relationship.FindAll."+amr.SrcResourceName+"."+amr.Name+"."+amr.DstResourceName, src.Id);
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
    panic("TODO");
}

func(amr *APIMountedRelationship) ResolveRecords(r *Request, lb RelationshipLinkRecords, record *Record, include *IncludeInstructions) (*ORelationship, []*Record) {
    rel := &ORelationship{
        IsSingle: lb.IsSingle(),
    };
    included := []*Record{};
    srcResource := r.API.GetResource(amr.SrcResourceName);
    dstResource := r.API.GetResource(amr.DstResourceName);
    r.API.Logger.Debugf("CALLING LINKRECORDS: %#v\n", record);
    records := lb.LinkRecords(r,srcResource,dstResource,record);
    for _, record := range records {
        rel.Data = append(rel.Data, record.GetResourceIdentifier());
    }
    return rel,included;
}
