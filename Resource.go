package jsonapi;

type RequestParams struct {
    Paginator Paginator
}

type Resource interface {
    GetFuture() Future
    //FindDefault(r *Request, rp RequestParams) Future
    //FindOne(r *Request, rp RequestParams, id FutureValue) Future
    //FindMany(r *Request, rp RequestParams, ids []FutureValue) Future
    //FindManyByField(r *Request, rp RequestParams, field, value FutureValue) Future
    // TODO: this iss necessary for optimizations but the backend
    // does not easily support this right now
    //FindManyByFieldWithManyValues(field string, []value string) ([]Ider, error)
    Delete(r *Request, id string)
    ParseJSON(r *Request, src *Record, raw []byte) (dst *Record)
    Create(r *Request, record *Record)
    Update(r *Request, record *Record)
}
