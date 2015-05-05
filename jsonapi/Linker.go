package jsonapi;

type Linker interface{
    Link(included *[]IderTyper) *OutputLinkageSet
}
