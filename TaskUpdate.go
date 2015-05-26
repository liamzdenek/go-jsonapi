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

    if err != nil {
        Reply(NewResponderBaseErrors(404, errors.New(fmt.Sprintf("Couldn't get record to be updated: %s", err))));
    }
    if ider == nil {
        Reply(NewResponderBaseErrors(404, errors.New(fmt.Sprintf("Record does not exist to be updated"))));
    }

    body, err := ioutil.ReadAll(r.Body);
    if err != nil {
        panic(NewResponderBaseErrors(400,errors.New(fmt.Sprintf("Body could not be parsed: %v\n", err))));
    }

    ider,id,rtype,linkages,err := resource.R.ParseJSON(s,ider,body);
    if err != nil {
        Reply(NewResponderBaseErrors(500, errors.New(fmt.Sprintf("ParseJSON threw error: %s", err))));
    }
    if(ider == nil) {
        Reply(NewResponderBaseErrors(500, errors.New("No error was thrown but ParseJSON did not return a valid object")));
    }
    if(rtype != nil && *rtype != resource_str) {
        Reply(NewResponderBaseErrors(409, errors.New(fmt.Sprintf("This is resource \"%s\" but the new object includes type:\"%s\"", resource_str, rtype))));
    }
    if(id != nil && *id != t.Id) {
        Reply(NewResponderBaseErrors(409, errors.New(fmt.Sprintf("The ID provided \"%s\" does not match the ID provided in the url \"%s\"", *id, t.Id))));
    }

    for _,linkage := range linkages.Linkages {
        fmt.Printf("LINKAGES: %#v\n", linkage);
        rel := a.RM.GetRelationship(resource_str, linkage.LinkName);
        err := rel.B.VerifyLinks(s, ider, linkage);
        if err != nil {
            
        }
    }

    err = resource.R.Update(s,t.Id,ider);
    if(err != nil) {
        Reply(NewResponderBaseErrors(500, errors.New(fmt.Sprintf("Could not update resource: %s", err))));
    }
    
    // TODO: the spec requires a 200 option with the requested resource if we modified it internally... no idea how to pull off that one
    Reply(NewResponderBase(202, nil));
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
