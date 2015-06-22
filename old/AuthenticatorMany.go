package jsonapi;

type AuthenticatorMany struct {
    Authenticators []Authenticator
}

func NewAuthenticatorMany(authenticators ...Authenticator) *AuthenticatorMany {
    return &AuthenticatorMany{
        Authenticators: authenticators,
    }
}

func(am *AuthenticatorMany) Authenticate(r *Request, permission, id string) {
    for _, authenticator := range am.Authenticators {
        authenticator.Authenticate(r,permission,id);
    }
}
