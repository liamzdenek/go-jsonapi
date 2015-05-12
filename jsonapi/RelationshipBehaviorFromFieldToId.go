package jsonapi;

import(
    "reflect"
    "strconv"
    "fmt"
);

type RelationshipBehaviorFromFieldToId struct {
    SrcFieldName string
}

func NewRelationshipBehaviorFromFieldToId(srcFieldName string) *RelationshipBehaviorFromFieldToId {
    return &RelationshipBehaviorFromFieldToId{
        SrcFieldName: srcFieldName,
    }
}

func(l *RelationshipBehaviorFromFieldToId) LinkId(srcR, dstR *ResourceManagerResource, src Ider) (ids []string) {
    v := reflect.Indirect(reflect.ValueOf(src)).FieldByName(l.SrcFieldName);
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
    fmt.Printf("Verify links\n");
    return nil;
}
func(l *RelationshipBehaviorFromFieldToId) PreCreate(ider Ider, linkages *OutputLinkage) error {
    fmt.Printf("Pre create\n");
    return nil;
}
func(l *RelationshipBehaviorFromFieldToId) PostCreate(ider Ider, linkages *OutputLinkage) error {
    fmt.Printf("Post create\n");
    return nil;
}
