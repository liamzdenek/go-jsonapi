package jsonapi;

type Record struct {
    // exposed fields
    Type string `json:"type"`
    Id string `json:"id"`
    Attributes interface{} `json:"attributes,omitempty"`
    //Links //TODO
    Relationships *ORelationships `json:"relationships,omitempty"`
    Meta OMeta `json:"meta,omitempty"`
    
    // internal fields for tracking
    ShouldInclude bool `json:"-"`
    relationshipsTask *TaskFindLinksByRecord `json:"-"`
}

func(r *Record) PrepareRelationships(req *Request, ii *IncludeInstructions) {
    if r.relationshipsTask == nil {
        r.relationshipsTask = NewTaskFindLinksByRecord(r,ii);
        req.Push(r.relationshipsTask);
    }
}

func(r *Record) GetRelationships() *TaskFindLinksByRecordResult {
    if r.relationshipsTask == nil {
        panic("Cannot call GetRelationships() before PrepareRelationships()")
    }
    return r.relationshipsTask.GetResult();
}

func(r *Record) GetResourceIdentifier() OResourceIdentifier {
    return OResourceIdentifier{
        Id: r.Id,
        Type: r.Type,
    }
}

func(r *Record) Denature() interface{} {
    return r.Attributes;
}
