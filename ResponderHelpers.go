package jsonapi;

import ("fmt";"errors";)

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
func NewResponderRecordCreate(resource_str string, rec *Record, createdStatus RecordCreatedStatus, err error) *ResponderBase {
    if(createdStatus & StatusCreated == 0) { // failure
        return NewResponderBaseErrors(500, err);
    } else {
        rb := NewResponderBase(201, nil);
        rb.CB = func(r *Request) {
            rb.PushHeader("Location", r.GetBaseURL()+resource_str+"/"+rec.Id);
        }
        return rb;
    }
}

func NewResponderResourceSuccessfullyDeleted() *ResponderBase {
    return NewResponderBase(204, nil);
}

func NewResponderForbidden(e error) *ResponderBase {
    return NewResponderBaseErrors(403, e);
}

func NewResponderUnimplemented(e error) *ResponderBase {
    return NewResponderBaseErrors(501, e);
}

func Unimplemented() *ResponderBase {
    err := "This endpoint is not implemented and will not be implemented."
    return NewResponderUnimplemented(errors.New(err));
}

func InsufficientPermissions() *ResponderBase {
    return NewResponderBaseErrors(403, errors.New("You do not have the required permission for this endpoint."));
}

func TODO() *ResponderBase {
    err := "This endpoint is not yet implemented."
    return NewResponderUnimplemented(errors.New(err));
}
