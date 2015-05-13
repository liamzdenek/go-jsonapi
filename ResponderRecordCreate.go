package jsonapi;

import("net/http");

type ResponderRecordCreate struct {
    Ider Ider
    Resource string
    Status RecordCreatedStatus
    Error error
}

func NewResponderRecordCreate(resource_str string, ider Ider, createdStatus RecordCreatedStatus, err error) *ResponderRecordCreate {
    return &ResponderRecordCreate{
        Resource: resource_str,
        Ider: ider,
        Status: createdStatus,
        Error: err,
    }
}

func (rrc *ResponderRecordCreate) Respond(a *API, w http.ResponseWriter, r *http.Request) error {
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
