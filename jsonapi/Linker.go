package jsonapi;

type Linker interface{
    Link(included *[]Record) *OutputLinkageSet
}
