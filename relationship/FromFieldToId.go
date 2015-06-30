package relationship;

import (
    . ".." // jsonapi
    "fmt"
    "reflect"
    "strconv"
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
    case *FutureResponseKindRecords:
        ids := []string{};
        for _, record := range t.Records {
            v := reflect.ValueOf(GetField(rel.SrcFieldName, record.Attributes));
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
