package jsonapi;

import("encoding/json";);

type Output struct { // responsible for the root node
    Data *OutputData `json:"data,omitempty"`
    //Links interface{} `json:"links,omitempty"`
    Included *OutputIncluded `json:"included,omitempty"`
    Errors []error `json:"errors,omitempty"`
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
    return json.Marshal(struct{
        Data *OutputData `json:"data,omitempty"`
        //Links interface{} `json:"links,omitempty"`
        Included *OutputIncluded `json:"included,omitempty"`
    }{
        Data: o.Data,
        Included: o.Included,
    });
}
