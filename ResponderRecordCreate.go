package jsonapi;

import("net/http");

type ResponderRecordCreate struct {
    Ider Ider
    Resource string
    Status RecordCreatedStatus
    Context Context
    Error error
}

func NewResponderRecordCreate(ctx Context, resource_str string, ider Ider, createdStatus RecordCreatedStatus, err error) *ResponderRecordCreate {
    return &ResponderRecordCreate{
        Resource: resource_str,
        Ider: ider,
        Status: createdStatus,
        Error: err,
        Context: ctx,
    }
}

func (rrc *ResponderRecordCreate) Respond(a *API, w http.ResponseWriter, r *http.Request) error {
    if rrc.Status & StatusFailed != 0 {
        err := rrc.Context.Failure();
        if err != nil {
            rrc.Error = err;
        }
    } else {
        err := rrc.Context.Success();
        if err != nil {
            rrc.Error = err;
        }
    }
    if(rrc.Error != nil) {
        res := NewResponderError(rrc.Error);
        res.Respond(a,w,r);
        return nil;
    }
    if rrc.Status & StatusCreated != 0 {
        w.WriteHeader(201) // 201 Created
        w.Header().Add("Location", a.GetBaseURL(r)+rrc.Resource+"/"+GetId(rrc.Ider));
    }
    //if rrc.Status & Status
    return nil;
}
