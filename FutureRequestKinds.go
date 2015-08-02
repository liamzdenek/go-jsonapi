package jsonapi

type FutureRequestKind interface{}

type FutureRequestKindFailure struct {
	Response *FutureResponse
}
type FutureRequestKindIdentity struct {
	Response FutureResponseKind
	Future
}
type FutureRequestKindFindByIds struct {
	Ids []string
}
type FutureRequestKindDeleteByIds struct {
	Ids []string
}

type Field struct {
	Field string
	Value string
}

type FutureRequestKindFindByAnyFields struct {
	Fields []Field
}
