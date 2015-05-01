package jsonapi;

type Linkage struct {
    Self string `json:"self"`
    Related string `json:"related"`
    Linkage []LinkageIdentifier `json:"linkage"`
}

type LinkageIdentifier struct {
    Type string `json:"type"`
    Id string `json:"id"`
}
