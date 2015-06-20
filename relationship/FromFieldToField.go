package relationship;

import(
    . ".."
    "errors"
    "fmt"
);

func init() {
    // safety check to make sure FromFieldToField is a Relationship and a RelationshipIder
    var t RelationshipLinkRecords = &FromFieldToField{};
    _ = t;
}

type FromFieldToField struct {
    SrcFieldName string
    DstResourceName string
    DstFieldName string
    Required RelationshipRequirement
    FromFieldToId *FromFieldToId
}

func NewFromFieldToField(dstResourceName, srcFieldName, dstFieldName string, required RelationshipRequirement) *FromFieldToField {
    return &FromFieldToField{
        SrcFieldName: srcFieldName,
        DstResourceName: dstResourceName,
        DstFieldName: dstFieldName,
        Required: required,
        FromFieldToId: NewFromFieldToId(dstResourceName, srcFieldName, required),
    }
}

func(l *FromFieldToField) IsSingle() (bool) { return false; }

func(l *FromFieldToField) PostMount(a *API) {
    if a.GetResource(l.DstResourceName) == nil {
        panic("FromFieldToId cannot be mounted to an API with a DstResourceName that does not exist");
    }
}

func(l *FromFieldToField) LinkRecords(r *Request, srcR *APIMountedResource, amr *APIMountedRelationship, src *Record) (dst []*Record) {
    a := r.API;
    ids := l.FromFieldToId.LinkIds(r,srcR,amr, src);
    dstR := a.GetResource(l.DstResourceName);
    //dstrmr := rmr.RM.GetResource(rmr.DstR);
    dst = []*Record{}
    for _, id := range ids {
        newdst, err := dstR.Resource.FindManyByField(r, RequestParams{}, l.DstFieldName, id.Id);
        if(err != nil) {
            a.Logger.Errorf("FromFieldToField got an error from FindManyByField for %s: %s", dstR.Name, err);
        }
        dst = append(dst, newdst...);
    }
    return dst;
}

func(l *FromFieldToField) VerifyLinks(r *Request, rec *Record, amr *APIMountedRelationship, linkages []OResourceIdentifier) error {
    isEmpty := linkages == nil || len(linkages) == 0;
    if(isEmpty && l.Required == Required) {
        return errors.New(fmt.Sprintf("Linkage '%s' is empty but is required",amr.Name));
    }
    return nil;
    //return l.FromFieldToId.VerifyLinks(s,ider,linkages);
}
func(l *FromFieldToField) PreSave(r *Request, rec *Record, amr *APIMountedRelationship, linkages []OResourceIdentifier) error {
    // Fetch the current links
    ii := NewIncludeInstructionsEmpty();
    ii.Push([]string{amr.Name});
    cur_links_task := NewTaskFindLinksByRecord(rec, ii);
    r.Push(cur_links_task);
    cur := cur_links_task.GetResult().Relationships.GetRelationshipByName(amr.Name)

    add, remove := GetRelationshipDifferences(cur.Data,linkages);
    if len(add) > 0 || len(remove) > 0 {
        panic(NewResponderUnimplemented(errors.New(
            fmt.Sprintf("FromFieldToField is a read-only relationship, and cannot be updated directly. If you wish to modify this relationship, you must either set the Source field on this record, '%s'. Or, you must set the value of the Target field, '%s', on the target Records of the target Resource, '%s'.", l.SrcFieldName, l.DstFieldName, l.DstResourceName),
        )));
    } else {
        panic("No differences");
    }
    return nil;
}
func(l *FromFieldToField) PostSave(r *Request, rec *Record, amr *APIMountedRelationship,linkages []OResourceIdentifier) error {
    return nil; 
}
