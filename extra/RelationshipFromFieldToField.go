package jsonapie;

import(
    . ".."
    "errors"
    "fmt"
);

func init() {
    // safety check to make sure RelationshipFromFieldToField is a Relationship and a RelationshipIder
    var t RelationshipLinkRecords = &RelationshipFromFieldToField{};
    _ = t;
}

type RelationshipFromFieldToField struct {
    SrcFieldName string
    DstResourceName string
    DstFieldName string
    Required RelationshipRequirement
    FromFieldToId *RelationshipFromFieldToId
}

func NewRelationshipFromFieldToField(dstResourceName, srcFieldName, dstFieldName string, required RelationshipRequirement) *RelationshipFromFieldToField {
    return &RelationshipFromFieldToField{
        SrcFieldName: srcFieldName,
        DstResourceName: dstResourceName,
        DstFieldName: dstFieldName,
        Required: required,
        FromFieldToId: NewRelationshipFromFieldToId(dstResourceName, srcFieldName, required),
    }
}

func(l *RelationshipFromFieldToField) IsSingle() (bool) { return false; }

func(l *RelationshipFromFieldToField) PostMount(a *API) {
    if a.GetResource(l.DstResourceName) == nil {
        panic("RelationshipFromFieldToId cannot be mounted to an API with a DstResourceName that does not exist");
    }
}

func(l *RelationshipFromFieldToField) LinkRecords(r *Request, srcR *APIMountedResource, amr *APIMountedRelationship, src *Record) (dst []*Record) {
    a := r.API;
    ids := l.FromFieldToId.LinkIds(r,srcR,amr, src);
    dstR := a.GetResource(l.DstResourceName);
    //dstrmr := rmr.RM.GetResource(rmr.DstR);
    dst = []*Record{}
    for _, id := range ids {
        newdst, err := dstR.Resource.FindManyByField(r, RequestParams{}, l.DstFieldName, id.Id);
        if(err != nil) {
            a.Logger.Errorf("RelationshipFromFieldToField got an error from FindManyByField for %s: %s", dstR.Name, err);
        }
        dst = append(dst, newdst...);
    }
    return dst;
}

func(l *RelationshipFromFieldToField) VerifyLinks(r *Request, rec *Record, amr *APIMountedRelationship, linkages []OResourceIdentifier) error {
    isEmpty := linkages == nil || len(linkages) == 0;
    if(isEmpty && l.Required == Required) {
        return errors.New(fmt.Sprintf("Linkage '%s' is empty but is required",amr.Name));
    }
    return nil;
    //return l.FromFieldToId.VerifyLinks(s,ider,linkages);
}
func(l *RelationshipFromFieldToField) PreSave(r *Request, rec *Record, amr *APIMountedRelationship, linkages []OResourceIdentifier) error {
    return nil; // no PreSave as we need Ider to be flushed to DB before we can use its ID
}
func(l *RelationshipFromFieldToField) PostSave(r *Request, rec *Record, amr *APIMountedRelationship,linkages []OResourceIdentifier) error {
    //id := rec.Id;
    //resource := a.GetResource(l.DstFieldName);
    //a := r.API;
    
    // Fetch the current links
    ii := NewIncludeInstructionsEmpty();
    ii.Push([]string{amr.Name});
    cur_links_task := NewTaskFindLinksByRecord(rec, ii);
    r.Push(cur_links_task);

    // remove the ones that shouldn't be there anymore
    cur_links := cur_links_task.GetResult().Relationships.GetRelationshipByName(amr.Name)
    OUTER: for _,cur_link := range cur_links.Data {
        for _,new_link := range linkages {
            if cur_link.Id == new_link.Id && cur_link.Type == new_link.Type {
                continue OUTER;
            }
        }
        // if we got to this point, the link exists in the current set but does not exist in the new set, and must be deleted
        panic("TODO: Asked to delete relationship");
    }
    
    // add ones that should be there now
    OUTER2: for _,new_link := range linkages {
        for _,cur_link := range cur_links.Data {
            if cur_link.Id == new_link.Id && cur_link.Type == new_link.Type {
                continue OUTER2;
            }
        }
        // if we got to this point, the link exists in the new set but does not exist in the old set, and must be added
        panic("TODO: asked to add relationship");
    }
    return nil;
}
