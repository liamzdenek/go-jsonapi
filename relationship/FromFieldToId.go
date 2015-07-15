package relationship;

import (
    . ".." // jsonapi
    "fmt"
    "reflect"
    "strconv"
    "strings"
);

type FromFieldToId struct {
    SrcFieldName string
}

func NewFromFieldToId(field string) *FromFieldToId {
    return &FromFieldToId{
        SrcFieldName: field,
    }
}

func(rel *FromFieldToId) IsSingle() bool { return true; }
func(rel *FromFieldToId) PostMount(a *API) { }
func(rel *FromFieldToId) Link(r *Request, src, dst *ExecutableFuture, input FutureResponseKind) (FutureRequestKind) {
    switch t := input.(type) {
    default:
        panic(fmt.Sprintf("FromFieldToId.Link does not support input of type %#T", input));
    case FutureResponseKindWithRecords:
        ids := []string{};
        for _, record := range t.GetRecords() {
            v := reflect.ValueOf(GetField(record.Attributes, rel.SrcFieldName));
            k := v.Kind()
            switch k { // TODO: fill this out
            case reflect.String:
                ids = append(ids, v.String());
            case reflect.Int:
                ids = append(ids, strconv.FormatInt(v.Int(), 10));
            default:
                panic("FromFieldToId does not know how to format the kind "+k.String());
            }
        }
        return &FutureRequestKindFindByIds{
            Ids: ids,
        }
    }
}
func(rel *FromFieldToId) PushBackRelationships(r *Request, src, dst *ExecutableFuture, srcrk, dstrk FutureResponseKind) {
    SimplePushBackRelationships(r,src,dst,srcrk,dstrk, map[string]string{
        rel.SrcFieldName: "Id",
    });
}

func SimplePushBackRelationships(r *Request, src, dst *ExecutableFuture, srcrk, dstrk FutureResponseKind, fieldmap map[string]string) {
    srcrkr, ok := srcrk.(FutureResponseKindWithRecords);
    if !ok {
        //panic("A");
        return;
    }
    dst.Request.API.Logger.Debugf("KINDBY FIELDS: %#v\n", dstrk);
    dstrkr, ok := dstrk.(*FutureResponseKindByFields);
    if !ok {
        //panic("B");
        return;
    }
    
    dst.Request.API.Logger.Debugf("PUSHING RELATIONSHIPS\n");
    for _, srcrecord := range srcrkr.GetRecords() {
        dst.Request.API.Logger.Debugf("PUSHING RELATIONSHIPS TO %s.%s\n", srcrecord.Type, srcrecord.Id);
        for field, records := range dstrkr.Records {
            if strings.ToLower(field.Field) == "id" {
                _, structfield := GetIdField(srcrecord.Attributes);
                field.Field = structfield.Name;
            } else {
                newfield, ok := fieldmap[field.Field];
                if !ok {
                    panic(fmt.Sprintf("Fieldmap did not contain a value for: %s\n %#v", field.Field, fieldmap));
                }
                field.Field = newfield;
            }
            dst.Request.API.Logger.Debugf("GOT FIELD: %#v\n", field);
            if srcrecord.HasFieldValue(field) {
                identifiers := GetResourceIdentifiers(records);
                newrel := &ORelationship{
                    IsSingle: dstrkr.IsSingle,
                    Data: identifiers,
                    RelationshipName: dst.Relationship.Name,
                    RelatedBase: dst.Request.GetBaseURL(),
                }
                dst.Request.API.Logger.Debugf("%s.%s IS GETTING %d NEW RELS\n", srcrecord.Type, srcrecord.Id, len(identifiers));
                dst.Request.API.Logger.Debugf("RELS: %#v\n", newrel);
                srcrecord.PushRelationship(newrel);
            }
        }
    }
}
