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
    r.API.Logger.Debugf("UPDATE got relationships: %#v\n", record.Relationships.Relationships);

    if record.Relationships == nil {
        record.Relationships = &ORelationships{};
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
        r.API.Logger.Debugf("CALLING PRE SAVE %s %s\n", resource_str, relationship.RelationshipName);
        err := rel.Relationship.PreSave(r, record, rel, relationship.Data);
        if err != nil {
            // TODO: A server MUST return 403 Forbidden in response to an unsupported request to update a resource or relationship. -- i don't know if this is the right behavior for this condition
            Reply(NewResponderBaseErrors(400, errors.New(fmt.Sprintf("Could not PreSave relationship %s: %s", relationship.RelationshipName, err))));
        }
    }

    err = resource.Resource.Update(r,record);
    if(err != nil) {
        Reply(NewResponderBaseErrors(500, errors.New(fmt.Sprintf("Could not update resource: %s", err))));
    }
    
    r.API.Logger.Infof("UPDATE WAS CALLED\n");
    for _,relationship := range record.Relationships.Relationships {
        rel := r.API.GetRelationship(resource_str, relationship.RelationshipName);
        err := rel.Relationship.PostSave(r, record, rel, relationship.Data);
        if err != nil {
            // TODO: A server MUST return 403 Forbidden in response to an unsupported request to update a resource or relationship. -- i don't know if this is the right behavior for this condition
            Reply(NewResponderBaseErrors(400, errors.New(fmt.Sprintf("Could not PostSave relationship %s: %s", relationship.RelationshipName, err))));
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
