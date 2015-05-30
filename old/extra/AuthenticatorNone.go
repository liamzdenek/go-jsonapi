package jsonapie;

import ("net/http";. "..");

type AuthenticatorNone struct{}

func init() {
    var a Authenticator = &AuthenticatorNone{}
    _ = a;
}

func NewAuthenticatorNone() *AuthenticatorNone {
    return &AuthenticatorNone{};
}

func (an *AuthenticatorNone) Authenticate(a *API, s Session, permission, id string, r *http.Request) {
    a.Logger.Printf("Authenticator request for: %s on ID: %s\n", permission, id);
}
