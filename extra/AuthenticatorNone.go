package jsonapie;

import ("net/http";"fmt");

type AuthenticatorNone struct{}

func NewAuthenticatorNone() *AuthenticatorNone {
    return &AuthenticatorNone{};
}

func (an *AuthenticatorNone) Authenticate(permission, id string, r *http.Request) {
    fmt.Printf("Authenticator request for: %s on ID: %s\n", permission, id);
}
