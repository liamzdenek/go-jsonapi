package jsonapi;

type Resource interface {
    FindOne(id string) (Ider, error)
    FindMany(ids []string) ([]Ider, error)
    FindManyByField(field, value string) ([]Ider, error)
    // TODO: this iss necessary for optimizations but the backend
    // does not easily support this right now
    //FindManyByFieldWithManyValues(field string, []value string) ([]Ider, error)
    Delete(id string) error
    ParseJSON(raw []byte) (ider Ider, id *string, rtype *string, links *OutputLinkageSet, err error)
    Create(ctx Context, resource_str string, ider Ider, id *string) (status RecordCreatedStatus, err error)
    //Update(resource_str, id string, ) error
}
