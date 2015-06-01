package jsonapi;

type APIMountedRelationship struct {
    SrcResourceName string
    DstResourceName string
    Name string
    Relationship
    Authenticator
}

func(amr *APIMountedRelationship) Resolve(r *Request, src *Record, shouldFetch bool, include *IncludeInstructions) (*ORelationship, []*Record) {
    return nil,nil
}
