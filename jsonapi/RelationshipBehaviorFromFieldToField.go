package jsonapi;

type RelationshipBehaviorFromFieldToField struct {
    SrcFieldName string
    DstFieldName string
    FromFieldToId *RelationshipBehaviorFromFieldToId
}

func NewRelationshipBehaviorFromFieldToField(srcFieldName, dstFieldName string) *RelationshipBehaviorFromFieldToField {
    return &RelationshipBehaviorFromFieldToField{
        SrcFieldName: srcFieldName,
        DstFieldName: dstFieldName,
        FromFieldToId: NewRelationshipBehaviorFromFieldToId(srcFieldName),
    }
}

func(l *RelationshipBehaviorFromFieldToField) Link(srcR, dstR *ResourceManagerResource, src Ider) (dst []Ider) {
    //ids := l.FromFieldToId.Link(src);
    //for _, id := range ids {
        
    //}
    return nil;
}
