package jsonapi;

type APIMountedRelationship struct {
    SrcResourceName string
    Name string
    Relationship
    Authenticator
}
