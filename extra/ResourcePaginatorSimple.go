package jsonapie;

import (
    . ".."
);

type ResourcePaginatorSimple struct{
    Parent Resource
    MaxPerPage int
}

func init() {
    // safety check to make sure ResourcePaginatorSimple is a Resource
    var t Resource;
    t = &ResourcePaginatorSimple{};
    _ = t;
}

func NewResourcePaginatorSimple(maxPerPage int, parent Resource) *ResourcePaginatorSimple {
    return &ResourcePaginatorSimple{
        Parent: parent,
        MaxPerPage: maxPerPage,
    }
}

func(rp *ResourcePaginatorSimple) FindOne(id string) (Ider, error) {
    return rp.Parent.FindOne(id);
}

func(rp *ResourcePaginatorSimple) FindMany(p *Paginator, ids []string) ([]Ider, error) {
    if p != nil {
        p.MaxPerPage = rp.MaxPerPage;
        p.LastPage = len(ids)/rp.MaxPerPage
    }
    return rp.Parent.FindMany(p, ids);
}

func(rp *ResourcePaginatorSimple) FindManyByField(field string, value string) ([]Ider, error) {
    return rp.Parent.FindManyByField(field, value);
}

func(rp *ResourcePaginatorSimple) Delete(id string) error {
    return rp.Parent.Delete(id);
}

func(rp *ResourcePaginatorSimple) ParseJSON(raw []byte) (Ider, *string, *string, *OutputLinkageSet, error) {
    return rp.Parent.ParseJSON(raw);
}

func(rp *ResourcePaginatorSimple) Create(ctx Context, resource_str string, ider Ider, id *string) (RecordCreatedStatus, error) {
    return rp.Parent.Create(ctx, resource_str, ider, id);
}
