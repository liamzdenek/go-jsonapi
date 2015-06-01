package jsonapi;

import ("encoding/json");

/*
OData is an interface used to represent the various forms of the primary data in the Output object. Currently, OutputData is an empty interface, so the use of this is purely syntactical. As of writing, the objects written in this file alongside OutputData are intended to be used as OutputData.
*/
type OData interface{}

/*
ORecords is a struct that satisfies OData, and represents a list of potentially many Records. A list of zero records is possible. Attempting to encode this object as JSON with a list of more than one Record when IsSingle = true will cause a panic. 
*/
type ORecords struct {
    IsSingle bool
    Records []*Record
}

func (o ORecords) MarshalJSON() ([]byte, error) {
    if o.IsSingle {
        if(len(o.Records) == 0) {
            return json.Marshal(nil);
        }
        return json.Marshal(o.Records[0]);
    }
    return json.Marshal(o.Records);
}

type ORelationships struct {
    Relationships []*ORelationship
    RelatedBase string
}

type ORelationship struct {
    IsSingle bool `json:"-"`
    //Links OLinks `json:"links,omitempty"`
    Data []OResourceIdentifier `json:"data"`
    Meta OMeta `json:"meta,omitempty"`
}

type OResourceIdentifier struct {
    Id string `json:"id"`
    Type string `json:"type"`
}
