package jsonapi;

import ("encoding/json";"errors");

type OutputLink struct { // data[i].links["linkname"].linkage[j]
    Type string `json:"type"`
    Id string `json:"id"`
}

type OutputLinkage struct { // data[i].links["linkname"] etc
    LinkName string
    Links []OutputLink `json:"linkage"`
}

func(o *OutputLinkage) UnmarshalJSON(data []byte) error {
    type A struct {
        Links []OutputLink `json:"linkage"`
    }
    type B struct {
        Link OutputLink `json:"linkage"`
    }
    a := &A{}
    err := json.Unmarshal(data, a);

    if err != nil {
        o.Links = a.Links;
        return nil;
    }

    b := &B{};
    err = json.Unmarshal(data, b);

    if err != nil {
        o.Links = []OutputLink{b.Link};
        return nil;
    }

    return errors.New("Linkages received were not valid");
}


type OutputLinkageSet struct { // data[i].links
    Linkages []*OutputLinkage
    RelatedBase string
    //Parent *OutputDatum
}

func(o *OutputLinkageSet) UnmarshalJSON(data []byte) error {
    res := map[string]*OutputLinkage{};
    err := json.Unmarshal(data, &res);
    if(err != nil) { return err; }
    return nil;
}

func(o *OutputLinkageSet) MarshalJSON() ([]byte,error) {
    /*if len(o.Linkages) == 0 {
        return json.Marshal(nil);
    }*/
    out := map[string]interface{}{};
    for _, linkage := range o.Linkages {
        out[linkage.LinkName] = struct{
            Links []OutputLink `json:"linkage"`
            Self string `json:"self"`
            Related string `json:"related"`
        }{
            Links:linkage.Links,
            Self:o.RelatedBase+"/links/"+linkage.LinkName,
            Related: o.RelatedBase+"/"+linkage.LinkName,
        };
    }
    //out["self"] = o.RelatedBase+"/links";
    return json.Marshal(out);
}
