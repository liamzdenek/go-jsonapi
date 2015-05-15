package jsonapi;

type LinkerStatic struct {
    Links *OutputLinkageSet
}

func NewLinkerStatic(links *OutputLinkageSet) *LinkerStatic {
    return &LinkerStatic{
        Links: links,
    }
}

func(ls *LinkerStatic) Link(included *[]Record) *OutputLinkageSet {
    return ls.Links
}
