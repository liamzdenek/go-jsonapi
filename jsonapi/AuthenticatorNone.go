package jsonapi;

import ("net/http");

type AuthenticatorNone struct{}

func NewAuthenticatorNone() *AuthenticatorNone {
    return &AuthenticatorNone{};
}

func (an *AuthenticatorNone) Authenticate(permission, id string, r *http.Request) {

}
