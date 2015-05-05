package jsonapi;

type Resource interface {
    FindOne(id string) (Ider, error)
    FindMany(ids []string) ([]Ider, error)
    FindManyByField(field string, value interface{}) ([]Ider, error)
}
