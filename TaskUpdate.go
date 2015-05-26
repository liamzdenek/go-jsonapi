package jsonapi;

import("net/http";"fmt";"io/ioutil";"errors");

type TaskUpdate struct {
    Resource, Id string
    Output chan chan bool
}

func NewTaskUpdate(resource, id string) *TaskUpdate {
    return &TaskUpdate{
        Resource: resource,
        Id: id,
        Output: make(chan chan bool),
    }
}

func(t *TaskUpdate) Work(a *API, s Session, tctx *TaskContext, r *http.Request) {
    resource_str := t.Resource;
    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(NewResponderErrorResourceDoesNotExist(resource_str));
    }

    resource.A.Authenticate(a,s,"resource.Update."+resource_str, "", r);

    ider, err := resource.R.FindOne(s,t.Id);

    body, err := ioutil.ReadAll(r.Body);
    if err != nil {
        panic(NewResponderError(errors.New(fmt.Sprintf("Body could not be parsed: %v\n", err))));
    }

    ider,id,rtype,linkages,err := resource.R.ParseJSON(s,ider,body);
    fmt.Printf("LINKAGES: %#v\n", linkages);
    if err != nil {
        Reply(NewResponderError(errors.New(fmt.Sprintf("ParseJSON threw error: %s", err))));
    }
    if(ider == nil) {
        Reply(NewResponderError(errors.New("No error was thrown but ParseJSON did not return a valid object")));
    }
    if(rtype != nil && *rtype != resource_str) {
        Reply(NewResponderError(errors.New(fmt.Sprintf("This is resource \"%s\" but the new object includes type:\"%s\"", resource_str, rtype))));
    }
    if(id != nil && *id != t.Id) {
        Reply(NewResponderError(errors.New(fmt.Sprintf("The ID provided \"%s\" does not match the ID provided in the url \"%s\"", id, t.Id))));
    }

    err = resource.R.Update(s,t.Id,ider);
    if(err != nil) {
        Reply(NewResponderError(errors.New(fmt.Sprintf("Could not update resource: %s", err))));
    }
    Reply("OK");
}

func(t *TaskUpdate) ResponseWorker(has_paniced bool) {
    go func() {
        for res := range t.Output {
            res <- true;
        }
    }();
}

func(t *TaskUpdate) Cleanup(a *API, r *http.Request) {

}

func(t *TaskUpdate) Wait() bool {
    r := make(chan bool);
    defer close(r);
    t.Output <- r;
    return <-r;
}
