package jsonapi;

type APIMountedRelationship struct {
    SrcResourceName string
    DstResourceName string
    Name string
    Relationship
    Authenticator
}

func(amr *APIMountedRelationship) Resolve(r *Request, src *Record, shouldFetch bool, include *IncludeInstructions) (*ORelationship, []*Record) {
    amr.Authenticator.Authenticate(r,"relationship.FindAll."+amr.SrcResourceName+"."+amr.Name+"."+amr.DstResourceName, src.Id);
    panic("TODO");
    return nil,nil
}
