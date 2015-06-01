package jsonapi;

/**
Output is the primary output structure used by this framework. It is responsible for representing the root node of every spec-compliant response this framework can generate.
*/
type Output struct {
    Data OData `json:"data,omitempty"`
    //Links *OutputLinks `json:"links, omitempty"`
    //Included *OutputIncluded `json:"included,omitempty"`
    Errors []OError `json:"errors,omitempty"`
    Meta OMeta `json:"meta,omitempty"`
}

func NewOutput() *Output {
    return &Output{};
}

type OMeta map[string]interface{}
