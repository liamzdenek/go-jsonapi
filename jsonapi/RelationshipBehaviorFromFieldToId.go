package jsonapi;

import(
    "reflect"
    "strconv"
    "fmt"
    "errors"
);

type RelationshipBehaviorFromFieldToId struct {
    SrcFieldName string
    Required RelationshipRequirement
}

func NewRelationshipBehaviorFromFieldToId(srcFieldName string, required RelationshipRequirement) *RelationshipBehaviorFromFieldToId {
    return &RelationshipBehaviorFromFieldToId{
        SrcFieldName: srcFieldName,
    }
}

func(l *RelationshipBehaviorFromFieldToId) LinkId(srcR, dstR *ResourceManagerResource, src Ider) (ids []string) {
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

func(l *RelationshipBehaviorFromFieldToId) VerifyLinks(ider Ider, linkages *OutputLinkage) error {
    isEmpty := linkages == nil || linkages.Links == nil || len(linkages.Links) == 0;
    fmt.Printf("LINKAGES: %#v\n", linkages);
    if(isEmpty && l.Required == Required) {
        return errors.New("Linkage is empty but is required");
    }
    if(!isEmpty && len(linkages.Links) != 1) {
        return errors.New("RelationshipBehaviorFromFieldToId requires exactly one link");
    }
    return nil;
}
func(l *RelationshipBehaviorFromFieldToId) PreCreate(ider Ider, linkages *OutputLinkage) error {
    str, err := strconv.Atoi(linkages.Links[0].Id);
    if err != nil {
        return err;
    }
    SetField(l.SrcFieldName, ider, str);
    fmt.Printf("Pre create\n");
    return nil;
}
func(l *RelationshipBehaviorFromFieldToId) PostCreate(ider Ider, linkages *OutputLinkage) error {
    fmt.Printf("Post create\n");
    return nil;
}
