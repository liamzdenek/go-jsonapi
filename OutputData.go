package jsonapi

import (
	"encoding/json"
)

/*
OData is an interface used to represent the various forms of the primary data in the Output object. Currently, OutputData is an empty interface, so the use of this is purely syntactical. As of writing, the objects written in this file alongside OutputData are intended to be used as OutputData.
*/
type OData interface{}

/*
ORecords is a struct that satisfies OData, and represents a list of potentially many Records. A list of zero records is possible. Attempting to encode this object as JSON with a list of more than one Record when IsSingle = true will cause a panic. 
*/
type ORecords struct {
	IsSingle bool
	Records  []*Record
}

func (o ORecords) MarshalJSON() ([]byte, error) {
	if o.IsSingle {
		if len(o.Records) == 0 {
			return json.Marshal(nil)
		}
		return json.Marshal(o.Records[0])
	}
	return json.Marshal(o.Records)
}

type ORelationships struct {
	Relationships []*ORelationship
}

func (o *ORelationships) MarshalJSON() ([]byte, error) {
	res := map[string]*ORelationship{}
	for _, relationship := range o.Relationships {
		res[relationship.RelationshipName] = relationship
	}
	return json.Marshal(res)
}

func (o *ORelationships) GetRelationshipByName(name string) *ORelationship {
	for _, relationship := range o.Relationships {
		if relationship.RelationshipName == name {
			return relationship
		}
	}
	return nil
}

type ORelationship struct {
	IsSingle bool `json:"-"`
	//Links OLinks `json:"links,omitempty"`
	Data                   []OResourceIdentifier `json:"data"`
	RelationshipWasFetched bool
	Meta                   OMeta  `json:"meta,omitempty"`
	RelatedBase            string `json:"-"`
	RelationshipName       string `json:"-"`
}

func (o *ORelationship) MarshalJSON() ([]byte, error) {
	links := map[string]string{
		"self":    o.RelatedBase + "/relationships/" + o.RelationshipName,
		"related": o.RelatedBase + "/" + o.RelationshipName,
	}
	if o.IsSingle && len(o.Data) == 0 {
		return json.Marshal(struct {
			Meta  OMeta       `json:"meta,omitempty"`
			Links interface{} `json:"links"`
		}{
			Meta:  o.Meta,
			Links: links,
		})
	} else if o.IsSingle {
		return json.Marshal(struct {
			Data  OResourceIdentifier `json:"data"`
			Meta  OMeta               `json:"meta,omitempty"`
			Links interface{}         `json:"links"`
		}{
			Data:  o.Data[0],
			Meta:  o.Meta,
			Links: links,
		})
	} else if o.RelationshipWasFetched {
		return json.Marshal(struct {
			Data  []OResourceIdentifier `json:"data"`
			Meta  OMeta                 `json:"meta,omitempty"`
			Links interface{}           `json:"links"`
		}{
			Data:  o.Data,
			Meta:  o.Meta,
			Links: links,
		})
	}
	return json.Marshal(struct {
		Data  []OResourceIdentifier `json:"data,omitempty"`
		Meta  OMeta                 `json:"meta,omitempty"`
		Links interface{}           `json:"links"`
	}{
		Data:  o.Data,
		Meta:  o.Meta,
		Links: links,
	})
}

type OResourceIdentifier struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

func NewResourceIdentifier(id, typ string) OResourceIdentifier {
	return OResourceIdentifier{
		Id:   id,
		Type: typ,
	}
}
