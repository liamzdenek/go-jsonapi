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

func(rp *ResourcePaginatorSimple) InitPaginator(p *Paginator) {
    if p != nil {
        p.MaxPerPage = rp.MaxPerPage;
        //p.LastPage = len(ids)/rp.MaxPerPage
    }
}

func(rp *ResourcePaginatorSimple) FinalizePaginator(p *Paginator, num_records int) {
    if p != nil {
        p.LastPage = num_records/rp.MaxPerPage
    }
}
func(rp *ResourcePaginatorSimple) FindDefault(a *API, s Session, p *Paginator) ([]Ider, error) {
    rp.InitPaginator(p);
    return rp.Parent.FindDefault(a,s,p);
}

func(rp *ResourcePaginatorSimple) FindOne(a *API, s Session, id string) (Ider, error) {
    return rp.Parent.FindOne(a, s, id);
}

func(rp *ResourcePaginatorSimple) FindMany(a *API, s Session, p *Paginator, ids []string) ([]Ider, error) {
    rp.InitPaginator(p);
    iders, err := rp.Parent.FindMany(a, s,p, ids);
    rp.FinalizePaginator(p, len(iders));
    return iders, err;
}

func(rp *ResourcePaginatorSimple) FindManyByField(a *API, s Session, field string, value string) ([]Ider, error) {
    return rp.Parent.FindManyByField(a,s,field, value);
}

func(rp *ResourcePaginatorSimple) Delete(a *API, s Session, id string) error {
    return rp.Parent.Delete(a,s,id);
}

func(rp *ResourcePaginatorSimple) ParseJSON(a *API, s Session, ider Ider, raw []byte) (Ider, *string, *string, *OutputLinkageSet, error) {
    return rp.Parent.ParseJSON(a, s, ider, raw);
}

func(rp *ResourcePaginatorSimple) Create(a *API, s Session, ider Ider, id *string) (RecordCreatedStatus, error) {
    return rp.Parent.Create(a,s, ider, id);
}

func(rp *ResourcePaginatorSimple) Update(a *API, s Session, id string, ider Ider) error {
    panic("NOT IMPLEMENTED");
}
