package jsonapi;

import ("net/http";"strconv");

type Paginator struct {
    CurPage int
    LastPage int
    MaxPerPage int
}

func NewPaginator(r *http.Request) *Paginator {
    page, err := strconv.Atoi(r.URL.Query().Get("page"));
    Check(err);
    return &Paginator{
        CurPage: page,
    }
}
