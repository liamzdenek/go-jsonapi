package jsonapi;

import("net/http";"fmt";"io/ioutil";"errors");

type TaskCreate struct {
    Resource, Id string
    Output chan chan bool
}

func NewTaskCreate(resource, id string) *TaskCreate {
    return &TaskCreate{
        Resource: resource,
        Id: id,
        Output: make(chan chan bool),
    }
}

func(t *TaskCreate) Work(a *API, s Session, tctx *TaskContext, r *http.Request) {
    resource_str := t.Resource;
    resource := a.RM.GetResource(resource_str);

    if(resource == nil) {
        panic(NewResponderErrorResourceDoesNotExist(resource_str));
    }

    resource.A.Authenticate(a,s,"resource.Create."+resource_str, "", r);

    body, err := ioutil.ReadAll(r.Body);
    if err != nil {
        panic(NewResponderError(errors.New(fmt.Sprintf("Body could not be parsed: %v\n", err))));
    }

    ider,id,rtype,linkages,err := resource.R.ParseJSON(s,nil,body);
    if err != nil {
        Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, errors.New(fmt.Sprintf("ParseJSON threw error: %s", err))));
    }

    if(ider == nil) {
        Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, errors.New("No error was thrown but ParseJSON did not return a valid object")));
    }
    if(rtype != nil && *rtype != resource_str) {
        Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, errors.New(fmt.Sprintf("This is resource \"%s\" but the new object includes type:\"%s\"", resource_str, rtype))));
    }
    if(id != nil && *id != "") {
        err = SetId(ider, *id);
        if(err != nil) {
            Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, errors.New(fmt.Sprintf("SetId failed:\"%s\"", err))));
        }
    }

    // first, we must check the permissions and verify that the
    // supplied linkages for each relationship is valid per the
    // rules of that relationship, eg, we don't want to let in
    // many linkages for a one to one relationship
    rels := a.RM.GetRelationshipsByResource(resource_str);
    for linkname, rel := range rels {
        linkage := linkages.GetLinkageByName(linkname);
        a.Logger.Printf("Linkage: %#v\n", linkage);
        if(rel == nil) {
            // user attempted to speify a relationship that does not exist
            panic("TODO: This");
        }
        rel.A.Authenticate(a,s,"relationship.Create."+rel.SrcR+"."+rel.Name+"."+rel.DstR, "", r);
        err := rel.B.VerifyLinks(a,s,ider,linkage);
        if err != nil {
            Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, err));
        }
    }
    // trigger the pre-creates so the linkages have a chance to modify
    // the id object before it's inserted
    for _,linkage := range linkages.Linkages {
        rel := a.RM.GetRelationship(resource_str, linkage.LinkName)
        err := rel.B.PreCreate(a,s,ider,linkage);
        if err != nil {
            Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, err));
        }
    }

    createdStatus, err := resource.R.Create(s,ider,id);
    if(err == nil && ider != nil && createdStatus & StatusCreated != 0) {
        for _,linkage := range linkages.Linkages {
            rel := a.RM.GetRelationship(resource_str, linkage.LinkName)
            err = rel.B.PostCreate(a,s,ider,linkage);
            if err != nil {
                Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, err));
            }
        }
    }
    Reply(NewResponderRecordCreate(resource_str, ider, createdStatus, err));
}

func(t *TaskCreate) ResponseWorker(has_paniced bool) {
    go func() {
        for res := range t.Output {
            res <- true;
        }
    }();
}

func(t *TaskCreate) Cleanup(a *API, r *http.Request) {

}

func(t *TaskCreate) Wait() bool {
    r := make(chan bool);
    defer close(r);
    t.Output <- r;
    return <-r;
}
