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

func(rp *ResourcePaginatorSimple) FindOne(a *API, s Session, id string) (Ider, error) {
    return rp.Parent.FindOne(a, s, id);
}

func(rp *ResourcePaginatorSimple) FindMany(a *API, s Session, p *Paginator, ids []string) ([]Ider, error) {
    if p != nil {
        p.MaxPerPage = rp.MaxPerPage;
        p.LastPage = len(ids)/rp.MaxPerPage
    }
    return rp.Parent.FindMany(a, s,p, ids);
}

func(rp *ResourcePaginatorSimple) FindManyByField(a *API, s Session, field string, value string) ([]Ider, error) {
    return rp.Parent.FindManyByField(a,s,field, value);
}

func(rp *ResourcePaginatorSimple) Delete(a *API, s Session, id string) error {
    return rp.Parent.Delete(a,s,id);
}

func(rp *ResourcePaginatorSimple) ParseJSON(a *API, s Session, raw []byte) (Ider, *string, *string, *OutputLinkageSet, error) {
    return rp.Parent.ParseJSON(a, s,raw);
}

func(rp *ResourcePaginatorSimple) Create(a *API, s Session, ider Ider, id *string) (RecordCreatedStatus, error) {
    return rp.Parent.Create(a,s, ider, id);
}
