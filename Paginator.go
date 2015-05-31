package jsonapi;

import ("strconv");

type Paginator struct {
    CurPage int
    LastPage int
    MaxPerPage int
}

func NewPaginator(r *Request) *Paginator {
    page, _ := strconv.Atoi(r.HttpRequest.URL.Query().Get("page"));
    // we do not care if the above line errors since we want the default value if it does, anyway
    return &Paginator{
        CurPage: page,
    }
}
