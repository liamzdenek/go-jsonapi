package jsonapi;

import("fmt";"net/http")

type WorkFindOne struct {
    Resource string
    Id string
    Output chan chan Ider
}

func NewWorkFindOne(id string, resource string) *WorkFindOne {
    return &WorkFindOne{
        Output: make(chan chan Ider),
        Id: id,
        Resource: resource,
    }
}

func(w *WorkFindOne) Work(a *API, r *http.Request) {
    resource := a.RM.GetResource(w.Resource);

    if(resource == nil) {
        panic(NewResponderErrorResourceDoesNotExist(w.Resource));
    }

    resource.A.Authenticate("resource.FindOne."+w.Resource, w.Id, r);

    data, err := resource.R.FindOne(w.Id);
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

func(w *WorkFindOne) Cleanup(a *API, r *http.Request) {
    fmt.Printf("INSIDE CLEANUP\n");
    close(w.Output);
}

func(w *WorkFindOne) GetResult() Ider {
    r := make(chan Ider);
    defer close(r);
    w.Output <- r;
    return <-r;
}
