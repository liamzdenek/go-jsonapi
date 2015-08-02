package relationship

import (
	. ".." // jsonapi
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type FromFieldToId struct {
	SrcFieldName    string
	DstResourceName string
	Required        RelationshipRequirement
}

func init() {
	// safety check to make sure FromFieldToId is a Relationship and a RelationshipLinkId
	var t RelationshipLinkIds = &FromFieldToId{}
	_ = t
}

func NewFromFieldToId(dstResourceName, srcFieldName string, required RelationshipRequirement) *FromFieldToId {
	return &FromFieldToId{
		SrcFieldName:    srcFieldName,
		DstResourceName: dstResourceName,
	}
}

func (l *FromFieldToId) IsSingle() bool { return true }

func (l *FromFieldToId) PostMount(a *API) {
	if a.GetResource(l.DstResourceName) == nil {
		panic("FromFieldToId cannot be mounted to an API with a DstResourceName that does not exist")
	}
}

func (l *FromFieldToId) LinkIds(r *Request, srcR *APIMountedResource, amr *APIMountedRelationship, src *Record) (ids []OResourceIdentifier) {
	r.API.Logger.Debugf("REFLECT FIELDS: %#v\n", src)
	v := reflect.ValueOf(GetField(l.SrcFieldName, src))
	k := v.Kind()
	switch k { // TODO: fill this out
	case reflect.String:
		ids = append(ids, NewResourceIdentifier(v.String(), l.DstResourceName))
	case reflect.Int:
		ids = append(ids, NewResourceIdentifier(strconv.FormatInt(v.Int(), 10), l.DstResourceName))
	default:
		panic("OneToOneLinkage does not support the kind " + k.String())
	}
	return ids
}

func (l *FromFieldToId) VerifyLinks(r *Request, record *Record, amr *APIMountedRelationship, linkages []OResourceIdentifier) error {
	a := r.API
	a.Logger.Infof("Verify links %#v\n", linkages)
	isEmpty := linkages == nil || len(linkages) == 0
	if isEmpty && l.Required == Required {
		return errors.New(fmt.Sprintf("Linkage '%s' is empty but is required", amr.Name))
	}
	if !isEmpty && len(linkages) != 1 {
		return errors.New("FromFieldToId requires exactly one link")
	}
	return nil
}
func (l *FromFieldToId) PreSave(r *Request, record *Record, amr *APIMountedRelationship, linkages []OResourceIdentifier) error {
	a := r.API
	a.Logger.Debugf("PreSave\n")
	if len(linkages) == 0 {
		return errors.New("FromFieldToId requires the relationship to be provided when modifying this relationship")
	}
	str, err := strconv.Atoi(linkages[0].Id)
	if err != nil {
		return err
	}
	SetField(l.SrcFieldName, record, str)
	return nil
}
func (l *FromFieldToId) PostSave(r *Request, record *Record, amr *APIMountedRelationship, linkages []OResourceIdentifier) error {
	a := r.API
	a.Logger.Debugf("Post create\n")
	return nil
}
