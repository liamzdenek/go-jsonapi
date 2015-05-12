package jsonapi;

type Linkage struct {
    Self string `json:"self,omitempty"`
    Related string `json:"related,omitempty"`
    Linkage []LinkageIdentifier `json:"linkage"`
}

type LinkageIdentifier struct {
    Type string `json:"type"`
    Id string `json:"id"`
}
