package jsonapi;

import( "fmt"; );

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

func(l *RelationshipBehaviorFromFieldToField) LinkIder(srcR, dstR *ResourceManagerResource, src Ider) (dst []Ider) {
    ids := l.FromFieldToId.LinkId(srcR, dstR, src);
    //dstrmr := rmr.RM.GetResource(rmr.DstR);
    dst = []Ider{}
    for _, id := range ids {
        newdst, err := dstR.R.FindManyByField(l.DstFieldName, id);
        if(err != nil) {
            fmt.Printf("RelationshipBehaviorFromFieldToField got an error from FindManyByField for %s: %s", dstR.Name, err);
        }
        dst = append(dst, newdst...);
    }
    return dst;
}

func(l *RelationshipBehaviorFromFieldToField) VerifyLinks(ider Ider, linkages *OutputLinkage) error {
    return nil;
}
func(l *RelationshipBehaviorFromFieldToField) PreCreate(ider Ider, linkages *OutputLinkage) error {
    fmt.Printf("Pre create\n");
    return nil;
}
func(l *RelationshipBehaviorFromFieldToField) PostCreate(ider Ider, linkages *OutputLinkage) error {
    fmt.Printf("Post create\n");
    return nil;
}
