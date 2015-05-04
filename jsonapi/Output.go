package jsonapi;

import("encoding/json";"net/http");

type Output struct { // responsible for the root node
    Data *OutputData `json:"data,omitempty"`
    //Links interface{} `json:"links,omitempty"`
    Included *OutputIncluded `json:"included,omitempty"`
    Errors []error `json:"errors,omitempty"`
    Request *http.Request `json:"-"`
}

func NewOutput(r *http.Request) *Output {
    return &Output{
        Data: &OutputData{},
        Included: &OutputIncluded{},
        Request: r,
    }
}

func (o Output) MarshalJSON() ([]byte, error) {
    // A document MUST contain either primary data or an array of error objects.
    if(len(o.Errors) > 0) {
        return json.Marshal(struct{
            Errors []error
        }{
            Errors: o.Errors,
        });
    }
    res := map[string]interface{}{};
    res["data"] = o.Data;
    if(o.Included.ShouldBeVisible()) {
        res["included"] = o.Included;
    }
    return json.Marshal(res);
}

