package jsonapie;

import (. "..");

type ContextFactorySimple struct {}

func NewContextFactorySimple() *ContextFactorySimple {
    return &ContextFactorySimple{};
}

func(cxm *ContextFactorySimple) NewContext() Context {
    return NewContextSimple();
}
