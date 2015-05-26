package jsonapie;

import(
    . ".."
);

func init() {
    // safety check to make sure RelationshipBehaviorFromFieldToField is a RelationshipBehavior and a RelationshipBehaviorIder
    var t RelationshipBehavior = &RelationshipBehaviorFromFieldToField{};
    _ = t;
    var t2 IderRelationshipBehavior = &RelationshipBehaviorFromFieldToField{};
    _ = t2;
}

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

func(l *RelationshipBehaviorFromFieldToField) LinkIder(s Session, srcR, dstR *ResourceManagerResource, src Ider) (dst []Ider) {
    a := s.GetData().API;
    ids := l.FromFieldToId.LinkId(s,srcR, dstR, src);
    //dstrmr := rmr.RM.GetResource(rmr.DstR);
    dst = []Ider{}
    for _, id := range ids {
        newdst, err := dstR.R.FindManyByField(s, RequestParams{}, l.DstFieldName, id);
        if(err != nil) {
            a.Logger.Printf("RelationshipBehaviorFromFieldToField got an error from FindManyByField for %s: %s", dstR.Name, err);
        }
        dst = append(dst, newdst...);
    }
    return dst;
}

func(l *RelationshipBehaviorFromFieldToField) VerifyLinks(s Session, ider Ider, linkages *OutputLinkage) error {
    panic("TODO");
    return l.FromFieldToId.VerifyLinks(s,ider,linkages);
}
func(l *RelationshipBehaviorFromFieldToField) PreSave(s Session, ider Ider, linkages *OutputLinkage) error {
    panic("TODO");
    return l.FromFieldToId.PreSave(s,ider,linkages);
}
func(l *RelationshipBehaviorFromFieldToField) PostSave(s Session, ider Ider, linkages *OutputLinkage) error {
    panic("TODO");
    return l.FromFieldToId.PostSave(s,ider,linkages);
}
