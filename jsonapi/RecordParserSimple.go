package jsonapi;

import("encoding/json";"errors")

type RecordParserSimple struct {
    Data *RecordParserSimpleData `json:"data"`
}

func NewRecordParserSimple(res interface{}) *RecordParserSimple {
    return &RecordParserSimple{
        Data: &RecordParserSimpleData{
            Output: res,
        },
    };
}

type RecordParserSimpleData struct {
    Output interface{} `json:"-"`
    Id *string
    Type string
}

func (rp *RecordParserSimpleData) UnmarshalJSON(data []byte) error {
    raw := map[string]interface{}{};
    err := json.Unmarshal(data, &raw);
    if(err != nil) {
        return err;
    }
    NatureObject(raw, rp.Output);
    if v, ok := raw["type"]; ok {
        rp.Type, ok = v.(string);
        if(!ok) {
            return errors.New("Type is required to be a string")
        }
    }
    if v, ok := raw["id"]; ok {
        nid, ok := v.(string);
        if(!ok) {
            return errors.New("Id is required to be a string");
        }
        rp.Id = &nid;
    }
    return nil;
}
