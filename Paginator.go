package jsonapi;

import ("net/http";"strconv");

type Paginator struct {
    CurPage int
    LastPage int
    MaxPerPage int
}

func NewPaginator(r *http.Request) *Paginator {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"));
    // we do not care if the above line errors since we want the default value if it does, anyway
    return &Paginator{
        CurPage: page,
    }
}
