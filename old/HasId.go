package jsonapi;

type HasId interface {
    GetId() string
    SetId(string) error
}
