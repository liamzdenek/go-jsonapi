package jsonapie;

import(
    . ".."
    "errors"
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
    Required RelationshipRequirement
    FromFieldToId *RelationshipBehaviorFromFieldToId
}

func NewRelationshipBehaviorFromFieldToField(srcFieldName, dstFieldName string, required RelationshipRequirement) *RelationshipBehaviorFromFieldToField {
    return &RelationshipBehaviorFromFieldToField{
        SrcFieldName: srcFieldName,
        DstFieldName: dstFieldName,
        Required: required,
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
    isEmpty := linkages == nil || linkages.Links == nil || len(linkages.Links) == 0;
    if(isEmpty && l.Required == Required) {
        return errors.New("Linkage is empty but is required");
    }
    return nil;
    //return l.FromFieldToId.VerifyLinks(s,ider,linkages);
}
func(l *RelationshipBehaviorFromFieldToField) PreSave(s Session, ider Ider, linkages *OutputLinkage) error {
    return nil; // no PreSave as we need Ider to be flushed to DB before we can use its ID
}
func(l *RelationshipBehaviorFromFieldToField) PostSave(s Session, ider Ider, linkages *OutputLinkage) error {
    panic("TODO");
    /*
    id := GetId(ider);
    a := s.GetData().API;
    wctx := s.GetData().TaskContext;
    resource := a.RM.GetResource(l.DstFieldName);
    
    // Fetch the current links
    ii := NewIncludeInstructionsEmpty();
    ii.Push([]string{linkages.LinkName});
    cur_links_task := NewWorkFindLinksByRecord(ider, ii);
    wctx.Push(cur_links_task);

    // remove the ones that shouldn't be there anymore
    cur_links := cur_links_task.GetResult().Links.GetLinkageByName(linkages.LinkName)
    OUTER: for _,cur_link := range cur_links.Links {
        for _,new_link := range linkages.Links {
            if cur_link.Id == new_link.Id && cur_link.Type == new_link.Type {
                continue OUTER;
            }
        }
        // if we got to this point, the link exists in the current set but does not exist in the new set, and must be deleted
        panic("TODO: Asked to delete linkage");
    }
    
    // add ones that should be there now
    OUTER: for _,new_link := range linkages.Links {
        for _,cur_link := range cur_links.Links {
            if cur_link.Id == new_link.Id && cur_link.Type == new_link.Type {
                continue OUTER;
            }
        }
        // if we got to this point, the link exists in the new set but does not exist in the old set, and must be added
        panic("TODO: asked to add linkage");
    }
    */
}
