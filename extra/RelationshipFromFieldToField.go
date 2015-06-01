package jsonapie;

import(
    . ".."
    "errors"
);

func init() {
    // safety check to make sure RelationshipFromFieldToField is a Relationship and a RelationshipIder
    var t RelationshipLinkRecords = &RelationshipFromFieldToField{};
    _ = t;
}

type RelationshipFromFieldToField struct {
    SrcFieldName string
    DstFieldName string
    Required RelationshipRequirement
    FromFieldToId *RelationshipFromFieldToId
}

func NewRelationshipFromFieldToField(srcFieldName, dstFieldName string, required RelationshipRequirement) *RelationshipFromFieldToField {
    return &RelationshipFromFieldToField{
        SrcFieldName: srcFieldName,
        DstFieldName: dstFieldName,
        Required: required,
        FromFieldToId: NewRelationshipFromFieldToId(srcFieldName, required),
    }
}

func(l *RelationshipFromFieldToField) IsSingle() (bool) { return false; }

func(l *RelationshipFromFieldToField) LinkRecords(r *Request, srcR, dstR *APIMountedResource, src *Record) (dst []*Record) {
    a := r.API;
    ids := l.FromFieldToId.LinkIds(r,srcR, dstR, src);
    //dstrmr := rmr.RM.GetResource(rmr.DstR);
    dst = []*Record{}
    for _, id := range ids {
        newdst, err := dstR.Resource.FindManyByField(r, RequestParams{}, l.DstFieldName, id);
        if(err != nil) {
            a.Logger.Errorf("RelationshipFromFieldToField got an error from FindManyByField for %s: %s", dstR.Name, err);
        }
        dst = append(dst, newdst...);
    }
    return dst;
}

func(l *RelationshipFromFieldToField) VerifyLinks(r *Request, rec *Record, linkages []*OResourceIdentifier) error {
    isEmpty := linkages == nil || len(linkages) == 0;
    if(isEmpty && l.Required == Required) {
        return errors.New("Linkage is empty but is required");
    }
    return nil;
    //return l.FromFieldToId.VerifyLinks(s,ider,linkages);
}
func(l *RelationshipFromFieldToField) PreSave(r *Request, rec *Record, linkages []*OResourceIdentifier) error {
    return nil; // no PreSave as we need Ider to be flushed to DB before we can use its ID
}
func(l *RelationshipFromFieldToField) PostSave(r *Request, rec *Record, linkages []*OResourceIdentifier) error {
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