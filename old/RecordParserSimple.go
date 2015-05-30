package jsonapi;

import("encoding/json";)

type RecordParserSimple struct {
    Data *RecordParserSimpleData `json:"data"`
}

func NewRecordParserSimple(res interface{}) *RecordParserSimple {
    return &RecordParserSimple{
        Data: &RecordParserSimpleData{
            Attributes: RecordParserSimpleAttributes{
                Output: res,
            },
        },
    };
}

func(rps *RecordParserSimple) Relationships() *OutputLinkageSet {
    return rps.Data.Relationships;
}

type RecordParserSimpleData struct {
    Attributes RecordParserSimpleAttributes `json:"attributes"`
    Id *string `json:"id"`
    Type string `json:"type"`
    Relationships *OutputLinkageSet `json:"relationships"`
}

type RecordParserSimpleAttributes struct {
    Output interface{} `json:"-"`
}

func (rp *RecordParserSimpleAttributes) UnmarshalJSON(data []byte) error {
    raw := map[string]interface{}{};

    err := json.Unmarshal(data, &raw);
    if(err != nil) {
        return err;
    }
    err = NatureObject(raw, rp.Output);
    return err;
}
