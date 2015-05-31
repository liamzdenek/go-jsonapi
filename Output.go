package jsonapi;

/**
Output is the primary output structure used by this framework. It is responsible for representing the root node of every spec-compliant response this framework can generate.
*/
type Output struct {
    //Data OutputData `json:"data,omitempty"`
    //Links *OutputLinks `json:"links, omitempty"`
    //Included *OutputIncluded `json:"included,omitempty"`
    Errors []OutputError `json:"errors,omitempty"`
    Meta interface{} `json:"meta,omitempty"`
}

func NewOutput() *Output {
    return &Output{};
}
