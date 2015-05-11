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

func(rps *RecordParserSimple) Linkages() *OutputLinkageSet {
    return rps.Data.Linkages;
}

type RecordParserSimpleData struct {
    Attributes RecordParserSimpleAttributes
    Id *string `json:"id"`
    Type string `json:"type"`
    Linkages *OutputLinkageSet `json:"links"`
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
    NatureObject(raw, rp.Output);
    return nil;
}
