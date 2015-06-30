package jsonapi;

type Record struct {
    // exposed fields
    Type string `json:"type"`
    Id string `json:"id"`
    Attributes RecordAttributes `json:"attributes,omitempty"`
    //Links //TODO
    Relationships *ORelationships `json:"relationships,omitempty"`
    Meta OMeta `json:"meta,omitempty"`

    // internal fields for tracking
    ShouldInclude bool `json:"-"`
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

type RecordAttributes interface{}
