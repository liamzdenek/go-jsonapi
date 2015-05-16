package jsonapi;

import("net/http";"strings";"errors");

type TaskDelete struct {
    Resource, Id string
    Result bool
    Output chan chan bool
}

func NewTaskDelete(resource, id string) *TaskDelete {
    return &TaskDelete{
        Resource: resource,
        Id: id,
        Output: make(chan chan bool),
    }
}

func(t *TaskDelete) Work(a *API, r *http.Request) {
    ids := strings.Split(t.Id,",");
    isSingle := len(ids) == 1;
    if(!isSingle) {
        panic(NewResponderError(errors.New("This request does not support more than one id")));
    }

    resource := a.RM.GetResource(t.Resource);

    if(resource == nil) {
        panic(NewResponderErrorResourceDoesNotExist(t.Resource));
    }

    resource.A.Authenticate("resource.Delete."+t.Resource, ids[0], r);

    err := resource.R.Delete(ids[0]);
    Check(err);
    t.Result = true;
    Reply(NewResponderResourceSuccessfullyDeleted());
}

func(t *TaskDelete) ResponseWorker(has_paniced bool) {
    go func() {
        for out := range t.Output {
            out <- t.Result;
        }
    }()
}

func(t *TaskDelete) Cleanup(a *API, r *http.Request) {
    close(t.Output);
}

func(t *TaskDelete) Wait() bool {
    r := make(chan bool);
    defer close(r);
    t.Output <- r;
    return <-r;
}
