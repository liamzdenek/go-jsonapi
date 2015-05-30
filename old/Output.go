package jsonapi;

import("encoding/json";"net/http";"strconv");

type Output struct { // responsible for the root node
    Data *OutputData `json:"data,omitempty"`
    Links *OutputPaginator `json:"links,omitempty"`
    Included *OutputIncluded `json:"included,omitempty"`
    Errors []OutputError `json:"errors,omitempty"`
    Meta interface{} `json:"meta,omitempty"`
}

type OutputPaginator struct {
    First string `json:"first,omitempty"`
    Prev string `json:"prev,omitempty"`
    Self string `json:"self,omitempty"`
    Next string `json:"next,omitempty"`
    Last string `json:"last,omitempty"`
}

func NewOutput(m interface{}) *Output {
    return &Output{
        Data: &OutputData{},
        Included: NewOutputIncluded(&[]Record{}),
        Meta: m,
    }
}

func (o *Output) Prepare() {
    if(o.Data != nil) {
        if(o.Data.Included == nil) {
            o.Data.Included = o.Included.Included;
        }
        o.Data.Prepare();
    }
}

func (o *Output) SetPaginator(r *http.Request, p *Paginator) {
    if p == nil || p.MaxPerPage == 0 {
        return;
    }
    q := r.URL.Query();
    proto := r.URL.Scheme;
    if proto == "" {
        proto = "http";
    }
    base := proto+"://"+r.Host+r.URL.Path
    l := &OutputPaginator{}

    q.Set("page", strconv.Itoa(p.CurPage));
    l.Self = base+"?"+q.Encode();
    q.Set("page", "0");
    l.First = base+"?"+q.Encode();

    if(p.CurPage > 0) {
        q.Set("page", strconv.Itoa(p.CurPage-1));
        l.Prev = base+"?"+q.Encode();
    }

    if(p.LastPage != 0) {
        q.Set("page", strconv.Itoa(p.LastPage));
        l.Last = base+"?"+q.Encode();
    }

    if(p.LastPage == 0 || p.CurPage < p.LastPage) {
        q.Set("page", strconv.Itoa(p.CurPage+1));
        l.Next = base+"?"+q.Encode();
    }
    o.Links = l;
}

func (o Output) MarshalJSON() ([]byte, error) {
    // A document MUST contain either primary data or an array of error objects.
    if(len(o.Errors) > 0) {
        //a.Logger.Printf("ERrors: %v\n", o.Errors);
        return json.Marshal(struct{
            Errors []OutputError `json:"errors"`
        }{
            Errors: o.Errors,
        });
    }
    res := map[string]interface{}{};
    res["data"] = o.Data;
    if(o.Meta != nil) {
        res["meta"] = o.Meta;
    }
    if(o.Included != nil && o.Included.ShouldBeVisible()) {
        res["included"] = o.Included;
    }
    if(o.Links != nil) {
        res["links"] = o.Links
    }
    return json.Marshal(res);
}

