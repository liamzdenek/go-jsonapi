package jsonapi;

type Resource interface {
    FindOne(id string) (Ider, error)
    FindMany(ids []string) ([]Ider, error)
    FindManyByField(field string, value string) ([]Ider, error)
    Delete(id string) error
    //Create(id string, ider Ider) error
}
