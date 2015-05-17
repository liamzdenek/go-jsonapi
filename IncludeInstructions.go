package jsonapi;

import ("net/http";"strings";
);

type IncludeInstructions struct {
    Children map[string]*IncludeInstructions
    Include []string
    Parent *IncludeInstructions
}

func NewIncludeInstructionsFromRequest(r *http.Request) *IncludeInstructions {
    e := NewIncludeInstructionsEmpty();
    e.Children["root"] = NewIncludeInstructions(r.URL.Query().Get("include"));
    return e;
}

func NewIncludeInstructionsEmpty() *IncludeInstructions {
    return &IncludeInstructions{
        Children: make(map[string]*IncludeInstructions),
    };
}

func NewIncludeInstructions(rawinst string) *IncludeInstructions {
    res := NewIncludeInstructionsEmpty();
    inst_strs := strings.Split(rawinst, ",")

    for _, inst_str := range inst_strs {
        inst_rels := strings.Split(inst_str,".");
        res.Push(inst_rels);
    }

    return res;
}

func(ii *IncludeInstructions) ShouldFetch(rel string) bool {
    if(ii.ShouldInclude(rel)) {
        return true;
    }
    _,ok := ii.Children[rel];
    //fmt.Printf("\nSHOULD FETCH %s: %s %s %#v\n\n", rel, ok, val, ii.Instructions);
    // TODO: do better
    return ok
}

func(ii *IncludeInstructions) ShouldInclude(inst string) bool {
    //fmt.Printf("Should include: %v %s\n", ii, inst);
    for _,included := range ii.Include {
        if(included == inst) {
            return true;
        }
    }
    return false;
}

func(ii *IncludeInstructions) Push(inst_rels []string) {
    if(len(inst_rels) == 0) {
        return;
    }
    if(len(inst_rels) == 1) {
        ii.Include = append(ii.Include, inst_rels[0]);
    } else {
        var child string;
        child, inst_rels = inst_rels[0], inst_rels[1:];
        if ii.Children[child] == nil {
            nii := NewIncludeInstructionsEmpty();
            nii.Parent = ii;
            ii.Children[child] = nii;
        }
        ii.Children[child].Push(inst_rels);
    }
}

func(ii *IncludeInstructions) GetChild(childname string) *IncludeInstructions {
    c := ii.Children[childname];
    if c == nil {
        c = NewIncludeInstructionsEmpty();
    }
    return c;
}
