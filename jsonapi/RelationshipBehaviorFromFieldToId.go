package jsonapi;

import(
    "reflect"
    "strconv"
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
