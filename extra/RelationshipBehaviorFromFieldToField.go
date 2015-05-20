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

func(l *RelationshipBehaviorFromFieldToField) LinkIder(a *API, s Session, srcR, dstR *ResourceManagerResource, src Ider) (dst []Ider) {
    ids := l.FromFieldToId.LinkId(a,s,srcR, dstR, src);
    //dstrmr := rmr.RM.GetResource(rmr.DstR);
    dst = []Ider{}
    for _, id := range ids {
        newdst, err := dstR.R.FindManyByField(a, s, l.DstFieldName, id);
        if(err != nil) {
            a.Logger.Printf("RelationshipBehaviorFromFieldToField got an error from FindManyByField for %s: %s", dstR.Name, err);
        }
        dst = append(dst, newdst...);
    }
    return dst;
}

func(l *RelationshipBehaviorFromFieldToField) VerifyLinks(a *API, s Session, ider Ider, linkages *OutputLinkage) error {
    return l.FromFieldToId.VerifyLinks(a,s,ider,linkages);
}
func(l *RelationshipBehaviorFromFieldToField) PreCreate(a *API, s Session, ider Ider, linkages *OutputLinkage) error {
    return l.FromFieldToId.PreCreate(a,s,ider,linkages);
}
func(l *RelationshipBehaviorFromFieldToField) PostCreate(a *API, s Session, ider Ider, linkages *OutputLinkage) error {
    return l.FromFieldToId.PostCreate(a,s,ider,linkages);
}
