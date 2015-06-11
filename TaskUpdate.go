package jsonapi;

import("fmt";"io/ioutil";"errors");

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

func(t *TaskUpdate) Work(r *Request) {
    resource_str := t.Resource;
    resource := r.API.GetResource(resource_str);

    if(resource == nil) {
        panic(NewResponderErrorResourceDoesNotExist(resource_str));
    }

    resource.Authenticator.Authenticate(r,"resource.Update."+resource_str, "");

    record, err := resource.Resource.FindOne(r,RequestParams{},t.Id);

    if err != nil {
        Reply(NewResponderBaseErrors(404, errors.New(fmt.Sprintf("Couldn't get record to be updated: %s", err))));
    }
    if record == nil {
        Reply(NewResponderBaseErrors(404, errors.New(fmt.Sprintf("Record does not exist to be updated"))));
    }

    body, err := ioutil.ReadAll(r.HttpRequest.Body);
    if err != nil {
        panic(NewResponderBaseErrors(400,errors.New(fmt.Sprintf("Body could not be parsed: %v\n", err))));
    }

    record,err = resource.Resource.ParseJSON(r,record,body);
    if err != nil {
        Reply(NewResponderBaseErrors(500, errors.New(fmt.Sprintf("ParseJSON threw error: %s", err))));
    }
    if(record == nil) {
        Reply(NewResponderBaseErrors(500, errors.New("No error was thrown but ParseJSON did not return a valid object")));
    }
    if(record.Type != "" && record.Type != resource_str) {
        Reply(NewResponderBaseErrors(409, errors.New(fmt.Sprintf("This is resource \"%s\" but the new object includes type:\"%s\"", resource_str, record.Type))));
    }
    if(record.Id != "" && record.Id != t.Id) {
        Reply(NewResponderBaseErrors(409, errors.New(fmt.Sprintf("The ID provided \"%s\" does not match the ID provided in the url \"%s\"", record.Id, t.Id))));
    }

    for _,relationship := range record.Relationships.Relationships {
        rel := r.API.GetRelationship(resource_str, relationship.RelationshipName);
        err := rel.Relationship.VerifyLinks(r, record, rel, relationship.Data);
        if err != nil {
            // TODO: A server MUST return 403 Forbidden in response to an unsupported request to update a resource or relationship. -- i don't know if this is the right behavior for this condition
            Reply(NewResponderBaseErrors(403, errors.New(fmt.Sprintf("Verification of new links for relationship %s failed: %s", relationship.RelationshipName, err))));
        }
    }

    for _,relationship := range record.Relationships.Relationships {
        rel := r.API.GetRelationship(resource_str, relationship.RelationshipName);
        a.Logger.Printf("CALLING PRE SAVE %s %s\n", resource_str, relationship.RelationshipName);
        err := rel.B.PreSave(s, record, linkage);
        if err != nil {
            // TODO: A server MUST return 403 Forbidden in response to an unsupported request to update a resource or relationship. -- i don't know if this is the right behavior for this condition
            Reply(NewResponderBaseErrors(400, errors.New(fmt.Sprintf("Could not PreSave relationship %s: %s", linkage.LinkName, err))));
        }
    }

    err = resource.R.Update(s,t.Id,record);
    if(err != nil) {
        Reply(NewResponderBaseErrors(500, errors.New(fmt.Sprintf("Could not update resource: %s", err))));
    }
    
    a.Logger.Printf("UPDATE WAS CALLED\n");
    for _,linkage := range linkages.Linkages {
        rel := a.RM.GetRelationship(resource_str, linkage.LinkName);
        err := rel.B.PostSave(s, record, linkage);
        if err != nil {
            // TODO: A server MUST return 403 Forbidden in response to an unsupported request to update a resource or relationship. -- i don't know if this is the right behavior for this condition
            Reply(NewResponderBaseErrors(400, errors.New(fmt.Sprintf("Could not PostSave relationship %s: %s", linkage.LinkName, err))));
        }
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

func(t *TaskUpdate) Cleanup(r *Request) {
    close(t.Output);
}

func(t *TaskUpdate) Wait() bool {
    r := make(chan bool);
    defer close(r);
    t.Output <- r;
    return <-r;
}
