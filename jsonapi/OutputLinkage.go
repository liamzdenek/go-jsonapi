package jsonapi;

import ("fmt";"encoding/json";);

type OutputLink struct { // data[i].links["linkname"].linkage[j]
    Type string `json:"type"`
    Id string `json:"id"`
}

type OutputLinkage struct { // data[i].links["linkname"] etc
    LinkName string
    Links []OutputLink
}

type OutputLinkageSet struct { // data[i].links
    Linkages []*OutputLinkage
    //Parent *OutputDatum
}

func(o *OutputLinkageSet) MarshalJSON() ([]byte,error) {
    out := map[string]interface{}{};
    for _, linkage := range o.Linkages {
        out[linkage.LinkName] = struct{
            Links []OutputLink `json:"linkage"`
            Self string `json:"self"`
            Related string `json:"related"`
        }{
            Links:linkage.Links,
        };
    }
    fmt.Printf("OutputLinkageSet.MarshalJSON: TODO: implement self field\n");
    out["self"] = "TODO";
    return json.Marshal(out);
}
