package jsonapi;

import ("encoding/json";"errors";"fmt");

// TODO: 2015-05-17 change to spec refers to this as a "Resource Identifier Object" -- should update this name properly
type OutputLink struct { // data[i].links["linkname"].linkage[j]
    Type string `json:"type"`
    Id string `json:"id"`
}

type OutputLinkage struct { // data[i].links["linkname"] etc
    LinkName string
    IsSingle bool
    Links []*OutputLink `json:"linkage"`
}

type OutputLinkageMany struct {
    Links []*OutputLink `json:"linkage"`
}
type OutputLinkageSingle struct {
    Link *OutputLink `json:"linkage"`
}

func(o *OutputLinkage) UnmarshalJSON(data []byte) error {
    a := &OutputLinkageSingle{}
    err := json.Unmarshal(data, a);

    if err == nil {
        o.Links = []*OutputLink{a.Link};
        return nil;
    }

    b := &OutputLinkageMany{};
    err = json.Unmarshal(data, b);

    if err == nil {
        o.Links = b.Links;
        return nil;
    }

    return errors.New("Linkages received were not valid");
}


type OutputLinkageSet struct { // data[i].links
    Linkages []*OutputLinkage
    RelatedBase string
    //Parent *OutputDatum
}

func(o *OutputLinkageSet) GetLinkageByName(name string) *OutputLinkage {
    for _, linkage := range o.Linkages {
        if(linkage.LinkName == name) {
            return linkage;
        }
    }
    return nil;
}

func(o *OutputLinkageSet) UnmarshalJSON(data []byte) error {
    fmt.Printf("\n\n\nUNMARSHAL JSON DATA: %s\n\n\n\n", data);
    type DataSingle struct{
        Data *OutputLink `json:"data"`
    }
    type DataMany struct {
        Data []*OutputLink `json:"data"`
    }
    links := map[string][]*OutputLink{};
    res_single := map[string]DataSingle{};
    err := json.Unmarshal(data, &res_single);
    if(err == nil) {
        for k, v := range res_single {
            links[k] = []*OutputLink{v.Data};
        }
    } else {
        res_many := map[string]DataMany{};
        err := json.Unmarshal(data, &res_many);
        if err != nil {
            return err;
        }
        for k,v := range res_many {
            links[k] = v.Data;
        }
    }
    fmt.Printf("LINKS: %#v\n", links);
    /*
    for linkname,r := range links {
        r.LinkName = linkname;
        o.Linkages = append(o.Linkages, r);
    }*/
    for k,vs := range links {
        o.Linkages = append(o.Linkages, &OutputLinkage{
            LinkName: k,
            Links: vs,
        })
    }
    return nil;
}

func(o *OutputLinkageSet) MarshalJSON() ([]byte,error) {
    /*if len(o.Linkages) == 0 {
        return json.Marshal(nil);
    }*/
    out := map[string]interface{}{};
    for _, linkage := range o.Linkages {
        if(linkage.IsSingle) {
            out[linkage.LinkName] = struct{
                Data *OutputLink `json:"data"`
                Links map[string]string `json:"links"`
            }{
                Data:linkage.Links[0],
                Links: map[string]string{
                    "self":o.RelatedBase+"/relationships/"+linkage.LinkName,
                    "related": o.RelatedBase+"/"+linkage.LinkName,
                },
            };
        } else {
            out[linkage.LinkName] = struct{
                Data []*OutputLink `json:"data"`
                Links map[string]string `json:"links"`
            }{
                Data:linkage.Links,
                Links: map[string]string{
                    "self":o.RelatedBase+"/relationships/"+linkage.LinkName,
                    "related": o.RelatedBase+"/"+linkage.LinkName,
                },
            };
        }
    }
    //out["self"] = o.RelatedBase+"/links";
    return json.Marshal(out);
}
