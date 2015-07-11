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

func(r *Record) PushRelationship(rel *ORelationship) {
    if(r.Relationships == nil) {
        r.Relationships = &ORelationships{};
    }
    for _, currel := range r.Relationships.Relationships {
        if currel.RelationshipName == rel.RelationshipName {
            currel.Data = append(currel.Data, rel.Data...);
            return;
        }
    }
    r.Relationships.Relationships = append(r.Relationships.Relationships, rel);
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

func GetResourceIdentifiers(records []*Record) (out []OResourceIdentifier) {
    for _, record := range records {
        out = append(out, record.GetResourceIdentifier());
    }
    return
}

type RecordAttributes interface{}
