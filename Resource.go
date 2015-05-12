package jsonapi;

type ResourceCreateVerifyFunc func(Ider, *OutputLinkageSet) error;

type Resource interface {
    FindOne(id string) (Ider, error)
    FindMany(ids []string) ([]Ider, error)
    FindManyByField(field, value string) ([]Ider, error)
    // TODO: this iss necessary for optimizations but the backend
    // does not easily support this right now
    //FindManyByFieldWithManyValues(field string, []value string) ([]Ider, error)
    Delete(id string) error
    Create(resource_str string, raw []byte, verify ResourceCreateVerifyFunc) (ider Ider, status RecordCreatedStatus, err error)
    //Update(resource_str, id string, ) error
}
