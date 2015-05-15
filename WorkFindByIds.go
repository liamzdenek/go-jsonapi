package jsonapi;

import("fmt";"net/http")

type WorkFindByIds struct {
    Resource string
    Ids []string
    Output chan chan []Ider
}

func NewWorkFindByIds(resource string, ids []string) *WorkFindByIds {
    return &WorkFindByIds{
        Output: make(chan chan []Ider),
        Ids: ids,
        Resource: resource,
    }
}

func(w *WorkFindByIds) Work(a *API, r *http.Request) {
    resource := a.RM.GetResource(w.Resource);

    if(resource == nil) {
        panic(NewResponderErrorResourceDoesNotExist(w.Resource));
    }

    // TODO: make this a loop over all the IDs
    for _, id := range w.Ids {
        resource.A.Authenticate("resource.FindOne."+w.Resource, id, r);
    }

    data := []Ider{}

    var err error;
    if(len(w.Ids) == 0) {
        panic("Ids must be longer than 0");
    } else if(len(w.Ids) == 1) {
        var ider Ider;
        ider, err = resource.R.FindOne(w.Ids[0]);
        data = []Ider{ider}
    } else {
        data, err = resource.R.FindMany(w.Ids);
    }
    if err != nil {
        // TODO: is this the right error?
        panic(NewResponderError(err));
    }
    go func() {
        for out := range w.Output {
            out <- data;
        }
    }();
}

func(w *WorkFindByIds) Cleanup(a *API, r *http.Request) {
    fmt.Printf("INSIDE CLEANUP\n");
    close(w.Output);
}

func(w *WorkFindByIds) GetResult() []Ider {
    r := make(chan []Ider);
    defer close(r);
    w.Output <- r;
    return <-r;
}
