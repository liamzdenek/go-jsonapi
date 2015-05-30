package jsonapi;

type Record struct {
    Type string `jsonapi:"type"`
    Id string `jsonapi:"id"`
    Attributes interface{} `jsonapi:"attributes,omitempty"`
    //Relationships
    //Links
    Meta interface{} `jsonapi:"meta,omitempty"`
}
