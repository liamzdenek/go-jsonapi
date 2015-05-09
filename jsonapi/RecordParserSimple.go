package jsonapi;

import("encoding/json")

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
}

func (rp *RecordParserSimpleData) UnmarshalJSON(data []byte) error {
    return json.Unmarshal(data, rp.Output);
}
