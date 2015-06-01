package jsonapie;

import(
    "reflect"
    "strconv"
    "errors"
    . ".." // jsonapi
);

type RelationshipBehaviorFromFieldToId struct {
    SrcFieldName string
    Required RelationshipRequirement
}

func init() {
    // safety check to make sure RelationshipBehaviorFromFieldToId is a RelationshipBehavior and a RelationshipBehaviorId
    var t RelationshipBehavior = &RelationshipBehaviorFromFieldToId{};
    _ = t;
    var t2 IdRelationshipBehavior = &RelationshipBehaviorFromFieldToId{};
    _ = t2;
}

func NewRelationshipBehaviorFromFieldToId(srcFieldName string, required RelationshipRequirement) *RelationshipBehaviorFromFieldToId {
    return &RelationshipBehaviorFromFieldToId{
        SrcFieldName: srcFieldName,
    }
}
func(l *RelationshipBehaviorFromFieldToId) IsSingle() (bool) { return true; }

func(l *RelationshipBehaviorFromFieldToId) LinkId(s Session, srcR, dstR *ResourceManagerResource, src Ider) (ids []string) {
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

func(l *RelationshipBehaviorFromFieldToId) VerifyLinks(s Session, ider Ider, linkages *OutputLinkage) error {
    a := s.GetData().API;
    a.Logger.Printf("Verify links %#v\n",linkages);
    isEmpty := linkages == nil || linkages.Links == nil || len(linkages.Links) == 0;
    if(isEmpty && l.Required == Required) {
        return errors.New("Linkage is empty but is required");
    }
    if(!isEmpty && len(linkages.Links) != 1) {
        return errors.New("RelationshipBehaviorFromFieldToId requires exactly one link");
    }
    return nil;
}
func(l *RelationshipBehaviorFromFieldToId) PreSave(s Session, ider Ider, linkages *OutputLinkage) error {
    a := s.GetData().API;
    a.Logger.Printf("PreSave\n");
    if(len(linkages.Links) == 0 || linkages.Links[0] == nil) {
        return errors.New("RelationshipBehaviorFromFieldToId requires the relationship to be provided when modifying this relationship");
    }
    str, err := strconv.Atoi(linkages.Links[0].Id);
    if err != nil {
        return err;
    }
    SetField(l.SrcFieldName, ider, str);
    return nil;
}
func(l *RelationshipBehaviorFromFieldToId) PostSave(s Session, ider Ider, linkages *OutputLinkage) error {
    a := s.GetData().API;
    a.Logger.Printf("Post create\n");
    return nil;
}