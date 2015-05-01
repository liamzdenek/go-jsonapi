package jsonapi;

// TODO: deprecate this
type TopLevel struct {
    Data interface{} `json:"data,omitempty"`
    Errors []error `json:"errors,omitempty"`
}
