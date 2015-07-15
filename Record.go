package jsonapi;

import "fmt" // for sprintf

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
    rel.RelationshipWasFetched = true;
    r.Relationships.Relationships = append(r.Relationships.Relationships, rel);
}

func(r *Record) GetResourceIdentifier() OResourceIdentifier {
    return OResourceIdentifier{
        Id: r.Id,
        Type: r.Type,
    }
}

func(r *Record) HasFieldValue(field Field) bool {
    val := GetField(r.Attributes, field.Field);
    // TODO: probably a better way to do this somehow
    if val != nil && (val == field.Value || fmt.Sprintf("%s", val) == field.Value || fmt.Sprintf("%d", val) == field.Value ){
        return true;
    }
    return false;
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
