package jsonapi;

//type RequestParams struct {
//    Paginator *Paginator
//}

type Resource interface {
    //FindDefault(s Session, rp RequestParams) ([]Ider, error)
    //FindOne(s Session, id string) (Ider, error)
    //FindMany(s Session, rp RequestParams, ids []string) ([]Ider, error)
    //FindManyByField(s Session, rp RequestParams, field, value string) ([]Ider, error)
    // TODO: this iss necessary for optimizations but the backend
    // does not easily support this right now
    //FindManyByFieldWithManyValues(field string, []value string) ([]Ider, error)
    //Delete(s Session, id string) error
    //ParseJSON(s Session, ider_src Ider, raw []byte) (ider Ider, id *string, rtype *string, links *OutputLinkageSet, err error)
    //Create(s Session, ider Ider, id *string) (status RecordCreatedStatus, err error)
    //Update(s Session, id string, ider Ider) error
}
