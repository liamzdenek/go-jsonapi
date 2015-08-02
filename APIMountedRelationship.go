package jsonapi

type APIMountedRelationship struct {
	SrcResourceName string
	DstResourceName string
	Name            string
	Relationship
	Authenticator
}
