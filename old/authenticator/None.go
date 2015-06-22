package authenticator;

import (. "..");

type None struct{}

func init() {
    var a Authenticator = &None{}
    _ = a;
}

func NewNone() *None {
    return &None{};
}

func (n *None) Authenticate(r *Request, permission, id string) {
    r.API.Logger.Infof("Authenticator request for: %s on ID: %s\n", permission, id);
}
