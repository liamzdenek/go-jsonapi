package jsonapi;

import ("fmt";"errors";"net/http")

func NewResponderErrorOperationNotSupported(desc string) *ResponderBase {
    return NewResponderBaseErrors(500, errors.New(fmt.Sprintf("The requested operation is not supported: %s\n", desc)));
}


func NewResponderErrorRelationshipDoesNotExist(relname string) *ResponderBase {
    return NewResponderBaseErrors(404, errors.New(fmt.Sprintf("The provided relationship \"%s\" does not exist.", relname)));
}

func NewResponderErrorResourceDoesNotExist(relname string) *ResponderBase {
    return NewResponderBaseErrors(404, errors.New(fmt.Sprintf("The provided relationship \"%s\" does not exist.", relname)));
}

// TODO: rip this out and replace it with multiple responder functions... this function should not be internally resonsible for determining success or failure
func NewResponderRecordCreate(resource_str string, ider Ider, createdStatus RecordCreatedStatus, err error) *ResponderBase {
    if(createdStatus & StatusCreated == 0) { // failure
        return NewResponderBaseErrors(500, err);
    } else {
        rb := NewResponderBase(201, nil);
        rb.CB = func(a *API, s Session, r *http.Request) {
            rb.PushHeader("Location", a.GetBaseURL(r)+resource_str+"/"+GetId(ider));
        }
        return rb;
    }
}

func NewResponderResourceSuccessfullyDeleted() *ResponderBase {
    return NewResponderBase(204, nil);
}
