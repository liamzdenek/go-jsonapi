package jsonapi;

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
