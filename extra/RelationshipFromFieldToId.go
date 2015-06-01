package jsonapie;

import(
    "reflect"
    "strconv"
    "errors"
    . ".." // jsonapi
);

type RelationshipFromFieldToId struct {
    SrcFieldName string
    Required RelationshipRequirement
}

func init() {
    // safety check to make sure RelationshipFromFieldToId is a Relationship and a RelationshipLinkId
    var t RelationshipLinkIds = &RelationshipFromFieldToId{};
    _ = t;
}

func NewRelationshipFromFieldToId(srcFieldName string, required RelationshipRequirement) *RelationshipFromFieldToId {
    return &RelationshipFromFieldToId{
        SrcFieldName: srcFieldName,
    }
}
func(l *RelationshipFromFieldToId) IsSingle() (bool) { return true; }

func(l *RelationshipFromFieldToId) LinkIds(r *Request, srcR, dstR *APIMountedResource, src *Record) (ids []string) {
    r.API.Logger.Debugf("REFLECT FIELDS: %#v\n", src);
    v := reflect.ValueOf(GetField(l.SrcFieldName, src));
    k := v.Kind()
    switch k { // TODO: fill this out
    case reflect.String:
        ids = append(ids, v.String());
    case reflect.Int:
        ids = append(ids, strconv.FormatInt(v.Int(), 10))
    default:
        panic("OneToOneLinkage does not support the kind "+k.String());
    }
    return ids;
}

func(l *RelationshipFromFieldToId) VerifyLinks(r *Request, record *Record, linkages []*OResourceIdentifier) error {
    a := r.API;
    a.Logger.Infof("Verify links %#v\n",linkages);
    isEmpty := linkages == nil || len(linkages) == 0;
    if(isEmpty && l.Required == Required) {
        return errors.New("Linkage is empty but is required");
    }
    if(!isEmpty && len(linkages) != 1) {
        return errors.New("RelationshipFromFieldToId requires exactly one link");
    }
    return nil;
}
func(l *RelationshipFromFieldToId) PreSave(r *Request, record *Record, linkages []*OResourceIdentifier) error {
    a := r.API;
    a.Logger.Debugf("PreSave\n");
    if(len(linkages) == 0 || linkages[0] == nil) {
        return errors.New("RelationshipFromFieldToId requires the relationship to be provided when modifying this relationship");
    }
    str, err := strconv.Atoi(linkages[0].Id);
    if err != nil {
        return err;
    }
    SetField(l.SrcFieldName, record, str);
    return nil;
}
func(l *RelationshipFromFieldToId) PostSave(r *Request, record *Record, linkages []*OResourceIdentifier) error {
    a := r.API;
    a.Logger.Debugf("Post create\n");
    return nil;
}
