package jsonapi;

import("encoding/json";"net/http";);

type Output struct { // responsible for the root node
    Data *OutputData `json:"data,omitempty"`
    Links *OutputPaginator `json:"links,omitempty"`
    Included *OutputIncluded `json:"included,omitempty"`
    Errors []error `json:"errors,omitempty"`
    Request *http.Request `json:"-"`
}

type OutputPaginator struct {
    First string `json:"first,omitempty"`
    Prev string `json:"prev,omitempty"`
    Self string `json:"self,omitempty"`
    Next string `json:"next,omitempty"`
    Last string `json:"last,omitempty"`
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
        //fmt.Printf("ERrors: %v\n", o.Errors);
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
    if(o.Links != nil) {
        res["links"] = o.Links
    }
    return json.Marshal(res);
}

