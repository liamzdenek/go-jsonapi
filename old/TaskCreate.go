package jsonapi

import (
	"errors"
	"fmt"
	"io/ioutil"
)

type TaskCreate struct {
	Resource, Id string
	Output       chan chan bool
}

func NewTaskCreate(resource, id string) *TaskCreate {
	return &TaskCreate{
		Resource: resource,
		Id:       id,
		Output:   make(chan chan bool),
	}
}

func (t *TaskCreate) Work(r *Request) {
	resource_str := t.Resource
	resource := r.API.GetResource(resource_str)

	if resource == nil {
		panic(NewResponderErrorResourceDoesNotExist(resource_str))
	}

	resource.Authenticator.Authenticate(r, "resource.Create."+resource_str, "")

	body, err := ioutil.ReadAll(r.HttpRequest.Body)
	if err != nil {
		panic(NewResponderBaseErrors(400, errors.New(fmt.Sprintf("Body could not be parsed: %v\n", err))))
	}

	record, err := resource.Resource.ParseJSON(r, nil, body)
	if err != nil {
		Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, errors.New(fmt.Sprintf("ParseJSON threw error: %s", err))))
	}

	if record == nil {
		Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, errors.New("No error was thrown but ParseJSON did not return a valid object")))
	}

	// first, we must check the permissions and verify that the
	// supplied linkages for each relationship is valid per the
	// rules of that relationship, eg, we don't want to let in
	// many linkages for a one to one relationship
	rels := r.API.GetRelationshipsByResource(resource_str)
	for relname, rel := range rels {
		relnew := record.Relationships.GetRelationshipByName(relname)
		if relnew == nil {
			relnew = &ORelationship{
				RelationshipName: relname,
			}
			record.Relationships.Relationships = append(record.Relationships.Relationships, relnew)
		}
		rel.Authenticator.Authenticate(r, "relationship.Create."+rel.SrcResourceName+"."+rel.Name, "")
		err := rel.Relationship.VerifyLinks(r, record, rel, relnew.Data)
		if err != nil {
			Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, err))
		}
	}
	// trigger the pre-creates so the linkages have a chance to modify
	// the record before it's inserted
	for relname, rel := range rels {
		relnew := record.Relationships.GetRelationshipByName(relname)
		err := rel.Relationship.PreSave(r, record, rel, relnew.Data)
		if err != nil {
			Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, err))
		}
	}

	createdStatus, err := resource.Resource.Create(r, record)
	if err == nil && createdStatus&StatusCreated != 0 {
		for relname, rel := range rels {
			relnew := record.Relationships.GetRelationshipByName(relname)
			err = rel.Relationship.PostSave(r, record, rel, relnew.Data)
			if err != nil {
				Reply(NewResponderRecordCreate(resource_str, nil, StatusFailed, err))
			}
		}
	}
	Reply(NewResponderRecordCreate(resource_str, record, createdStatus, err))
}

func (t *TaskCreate) ResponseWorker(has_paniced bool) {
	go func() {
		for res := range t.Output {
			res <- true
		}
	}()
}

func (t *TaskCreate) Cleanup(r *Request) {
	close(t.Output)
}

func (t *TaskCreate) Wait() bool {
	r := make(chan bool)
	defer close(r)
	t.Output <- r
	return <-r
}
