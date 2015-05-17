package jsonapie;

import(
    "fmt";
    . ".."
);

type RelationshipBehaviorFromFieldToField struct {
    SrcFieldName string
    DstFieldName string
    FromFieldToId *RelationshipBehaviorFromFieldToId
}

func NewRelationshipBehaviorFromFieldToField(srcFieldName, dstFieldName string, required RelationshipRequirement) *RelationshipBehaviorFromFieldToField {
    return &RelationshipBehaviorFromFieldToField{
        SrcFieldName: srcFieldName,
        DstFieldName: dstFieldName,
        FromFieldToId: NewRelationshipBehaviorFromFieldToId(srcFieldName, required),
    }
}

func(l *RelationshipBehaviorFromFieldToField) IsSingle() (bool) { return false; }
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
    return l.FromFieldToId.VerifyLinks(ider,linkages);
}
func(l *RelationshipBehaviorFromFieldToField) PreCreate(ider Ider, linkages *OutputLinkage) error {
    return l.FromFieldToId.PreCreate(ider,linkages);
}
func(l *RelationshipBehaviorFromFieldToField) PostCreate(ider Ider, linkages *OutputLinkage) error {
    return l.FromFieldToId.PostCreate(ider,linkages);
}
