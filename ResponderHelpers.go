package jsonapi;

import ("fmt";"errors";"net/http")

func NewResponderErrorsBase(code int, es ...error) *ResponderBase {
    // TODO: replace this with code to make NewResponderErrors accept a list of OutputError to begin with
    oes := []OutputError{}
    for _,e := range es {
        oes = append(oes, OutputError{
            Title: e.Error(),
        });
    }
    o := NewOutput(nil);
    o.Errors = oes;
    return NewResponderBase(code, o);
}

func NewResponderError(e error) *ResponderBase {
    return NewResponderErrors([]error{e});
}

// TODO: deprecate this. all ResponderErrors should have specific funcs and this generic function encourages laziness when an error happens
func NewResponderErrors(es []error) *ResponderBase {
    return NewResponderErrorsBase(500, es...);
}

func NewResponderErrorOperationNotSupported(desc string) *ResponderBase {
    return NewResponderErrorsBase(500, errors.New(fmt.Sprintf("The requested operation is not supported: %s\n", desc)));
}


func NewResponderErrorRelationshipDoesNotExist(relname string) *ResponderBase {
    return NewResponderErrorsBase(404, errors.New(fmt.Sprintf("The provided relationship \"%s\" does not exist.", relname)));
}

func NewResponderErrorResourceDoesNotExist(relname string) *ResponderBase {
    return NewResponderErrorsBase(404, errors.New(fmt.Sprintf("The provided relationship \"%s\" does not exist.", relname)));
}

// TODO: rip this out and replace it with multiple responder functions... this function should not be internally resonsible for determining success or failure
func NewResponderRecordCreate(resource_str string, ider Ider, createdStatus RecordCreatedStatus, err error) *ResponderBase {
    if(createdStatus & StatusCreated == 0) { // failure
        return NewResponderErrorsBase(500, err);
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
