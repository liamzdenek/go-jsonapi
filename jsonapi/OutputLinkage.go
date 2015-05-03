package jsonapi;

import ("encoding/json";);

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
    RelatedBase string
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
            Self:o.RelatedBase+"/links/"+linkage.LinkName,
            Related: o.RelatedBase+"/"+linkage.LinkName,
        };
    }
    out["self"] = o.RelatedBase+"/links";
    return json.Marshal(out);
}
