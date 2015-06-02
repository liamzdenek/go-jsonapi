package jsonapi;

import("strings";"errors");

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

func(t *TaskDelete) Work(r *Request) {
    ids := strings.Split(t.Id,",");
    isSingle := len(ids) == 1;
    if(!isSingle) {
        panic(NewResponderBaseErrors(400,errors.New("This request does not support more than one id")));
    }

    resource := r.API.GetResource(t.Resource);

    if(resource == nil) {
        panic(NewResponderErrorResourceDoesNotExist(t.Resource));
    }

    resource.Authenticator.Authenticate(r,"resource.Delete."+t.Resource, ids[0]);

    err := resource.Resource.Delete(r, ids[0]);
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

func(t *TaskDelete) Cleanup(r *Request) {
    close(t.Output);
}

func(t *TaskDelete) Wait() bool {
    r := make(chan bool);
    defer close(r);
    t.Output <- r;
    return <-r;
}
