package jsonapi;

import("encoding/json";"net/http";"fmt");

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
        Included: NewOutputIncluded(&[]Record{}),
        Request: r,
    }
}

func (o *Output) Prepare() {
    if(o.Data.Included == nil) {
        o.Data.Included = o.Included.Included;
    }
    o.Data.Prepare();
}

func (o Output) MarshalJSON() ([]byte, error) {
    // A document MUST contain either primary data or an array of error objects.
    if(len(o.Errors) > 0) {
        fmt.Printf("ERrors: %v\n", o.Errors);
        es := []string{};
        for _,e := range o.Errors {
            es = append(es, e.Error());
        }
        return json.Marshal(struct{
            Errors []string `json:"errors"`
        }{
            Errors: es,
        });
    }
    res := map[string]interface{}{};
    res["data"] = o.Data;
    if(o.Included != nil && o.Included.ShouldBeVisible()) {
        res["included"] = o.Included;
    }
    return json.Marshal(res);
}

