package jsonapi;

type RequestParams struct {
    Paginator *Paginator

}

type Resource interface {
    FindDefault(r *Request, rp RequestParams) ([]*Record, error)
    FindOne(r *Request, rp RequestParams, id string) (*Record, error)
    FindMany(r *Request, rp RequestParams, ids []string) ([]*Record, error)
    FindManyByField(r *Request, rp RequestParams, field, value string) ([]*Record, error)
    // TODO: this iss necessary for optimizations but the backend
    // does not easily support this right now
    //FindManyByFieldWithManyValues(field string, []value string) ([]Ider, error)
    Delete(r *Request, id string) error
    ParseJSON(r *Request, src *Record, raw []byte) (dst *Record, err error)
    //Create(r *Request, record *Record) (status RecordCreatedStatus, err error)
    Update(r *Request, record *Record) error
}
