package jsonapi;

type OutputError struct {
    Id string `json:"id,omitempty"`
    Href string `json:"href,omitempty"`
    Status string `json:"status,omitempty"`
    Code string `json:"code,omitempty"`
    Title string `json:"title,omitempty"`
    Detail string `json:"detail,omitempty"`
    Source *OutputErrorSource `json:"source,omitempty"`
    Meta interface{} `json:"meta,omitempty"`
}

type OutputErrorSource struct {
    Pointer string `json:"pointer,omitempty"`
    Parameter string `json:"parameter,omitempty"`
}
