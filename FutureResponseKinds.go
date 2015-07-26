package jsonapi

type FutureResponseKind interface {
}
type FutureResponseKindWithRecords interface {
	GetIsSingle() bool
	GetRecords() []*Record
}

type FutureResponseKindByFields struct {
	IsSingle bool
	Records  map[Field][]*Record
}

func (frkbf *FutureResponseKindByFields) GetIsSingle() bool {
	return false
}
func (frkbf *FutureResponseKindByFields) GetRecords() []*Record {
	rec := []*Record{}
	for _, recs := range frkbf.Records {
		rec = append(rec, recs...)
	}
	return rec
}

type FutureResponseKindDeleted struct {
}
