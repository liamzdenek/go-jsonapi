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
func(rp *ResourcePaginatorSimple) FindDefault(s Session, rparams RequestParams) ([]Ider, error) {
    rp.InitPaginator(rparams.Paginator);
    return rp.Parent.FindDefault(s,rparams);
}

func(rp *ResourcePaginatorSimple) FindOne(s Session, id string) (Ider, error) {
    return rp.Parent.FindOne(s, id);
}

func(rp *ResourcePaginatorSimple) FindMany(s Session, rparams RequestParams, ids []string) ([]Ider, error) {
    rp.InitPaginator(rparams.Paginator);
    iders, err := rp.Parent.FindMany(s,rparams,ids);
    rp.FinalizePaginator(rparams.Paginator, len(iders));
    return iders, err;
}

func(rp *ResourcePaginatorSimple) FindManyByField(s Session, rparams RequestParams, field string, value string) ([]Ider, error) {
    rp.InitPaginator(rparams.Paginator);
    iders, err := rp.Parent.FindManyByField(s,rparams,field, value);
    rp.FinalizePaginator(rparams.Paginator, len(iders));
    return iders, err;
}

func(rp *ResourcePaginatorSimple) Delete(s Session, id string) error {
    return rp.Parent.Delete(s,id);
}

func(rp *ResourcePaginatorSimple) ParseJSON(s Session, ider Ider, raw []byte) (Ider, *string, *string, *OutputLinkageSet, error) {
    return rp.Parent.ParseJSON(s, ider, raw);
}

func(rp *ResourcePaginatorSimple) Create(s Session, ider Ider, id *string) (RecordCreatedStatus, error) {
    return rp.Parent.Create(s, ider, id);
}

func(rp *ResourcePaginatorSimple) Update(s Session, id string, ider Ider) error {
    panic("NOT IMPLEMENTED");
}
