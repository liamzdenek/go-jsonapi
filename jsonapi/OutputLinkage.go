package jsonapi;

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
}
