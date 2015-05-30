package jsonapie;

import (. "..");

type SessionFactorySimple struct {}

func NewSessionFactorySimple() *SessionFactorySimple {
    return &SessionFactorySimple{};
}

func(cxm *SessionFactorySimple) NewSession() Session {
    return NewSessionSimple();
}
