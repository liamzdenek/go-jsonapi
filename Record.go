package jsonapi;

type Record struct {
    Type string `jsonapi:"type"`
    Id string `jsonapi:"id"`
    Attributes interface{} `jsonapi:"attributes,omitempty"`
    Meta OMeta `jsonapi:"meta,omitempty"`
    //Links
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
