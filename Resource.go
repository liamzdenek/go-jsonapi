package jsonapi;

type Resource interface {
    FindOne(a *API, s Session, id string) (Ider, error)
    FindMany(a *API, s Session, page *Paginator, ids []string) ([]Ider, error)
    FindManyByField(a *API, s Session, field, value string) ([]Ider, error)
    // TODO: this iss necessary for optimizations but the backend
    // does not easily support this right now
    //FindManyByFieldWithManyValues(field string, []value string) ([]Ider, error)
    Delete(a *API, s Session, id string) error
    ParseJSON(a *API, s Session, raw []byte) (ider Ider, id *string, rtype *string, links *OutputLinkageSet, err error)
    Create(a *API, s Session, ider Ider, id *string) (status RecordCreatedStatus, err error)
    //Update(resource_str, id string, ) error
}
