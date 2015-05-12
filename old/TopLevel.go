package jsonapi;

// TODO: deprecate this
type TopLevel struct {
    Data interface{} `json:"data,omitempty"`
    Links interface{} `json:"links,omitempty"`
    Included interface{} `json:"included,omitempty"`
    Errors []error `json:"errors,omitempty"`
}
