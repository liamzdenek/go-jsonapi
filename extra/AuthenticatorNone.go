package jsonapie;

import (. "..");

type AuthenticatorNone struct{}

func init() {
    var a Authenticator = &AuthenticatorNone{}
    _ = a;
}

func NewAuthenticatorNone() *AuthenticatorNone {
    return &AuthenticatorNone{};
}

func (an *AuthenticatorNone) Authenticate(r *Request, permission, id string) {
    //r.API.Logger.Printf("Authenticator request for: %s on ID: %s\n", permission, id);
}
