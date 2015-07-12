package relationship;

import (
    . ".." // jsonapi
    "fmt"
    "reflect"
    "strconv"
);

type FromFieldToField struct {
    SrcFieldName, DstFieldName string
}

func NewFromFieldToField(srcField, dstField string) *FromFieldToField {
    return &FromFieldToField{
        SrcFieldName: srcField,
        DstFieldName: dstField,
    }
}


func(rel *FromFieldToField) IsSingle() bool { return false; }
func(rel *FromFieldToField) PostMount(a *API) { }
func(rel *FromFieldToField) Link(r *Request, src, dst *ExecutableFuture, input FutureResponseKind) (FutureRequestKind) {
    switch t := input.(type) {
    default:
        panic(fmt.Sprintf("FromFieldToId.Link does not support input of type %#T", input));
    case FutureResponseKindWithRecords:
        fields := []Field{};
        for _, record := range t.GetRecords() {
            v := reflect.ValueOf(GetField(record.Attributes, rel.SrcFieldName));
            k := v.Kind()
            switch k { // TODO: fill this out
            case reflect.String:
                fields = append(fields, Field{
                    Field:rel.DstFieldName,
                    Value:v.String(),
                });
            case reflect.Int:
                fields = append(fields, Field{
                    Field:rel.DstFieldName,
                    Value:strconv.FormatInt(v.Int(), 10),
                });
            default:
                panic("FromFieldToId does not know how to format the kind "+k.String());
            }
        }
        return &FutureRequestKindFindByAnyFields{
            Fields: fields,
        }
    }
}
