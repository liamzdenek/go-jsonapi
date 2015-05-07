package jsonapi;

import ("net/http";"strings");

type IncludeInstructions struct {
    Instructions map[string]bool
    Children map[string]*IncludeInstructions
}

func NewIncludeInstructionsFromRequest(r *http.Request) *IncludeInstructions {
    return NewIncludeInstructions(r.URL.Query().Get("include"));
}

func NewIncludeInstructionsEmpty() *IncludeInstructions {
    return &IncludeInstructions{
        Children: make(map[string]*IncludeInstructions),
        Instructions: make(map[string]bool),
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

func(ii *IncludeInstructions) Handling(inst string) bool {
    res := !ii.Instructions[inst];
    ii.Instructions[inst] = true;
    return res;
}

func(ii *IncludeInstructions) Push(inst_rels []string) {
    if(len(inst_rels) == 0) {
        return;
    }
    if(len(inst_rels) == 1) {
        if(len(inst_rels[0]) > 0) {
            ii.Instructions[inst_rels[0]] = false;
        }
    } else {
        var child string;
        child, inst_rels = inst_rels[0], inst_rels[1:];
        if ii.Children[child] == nil {
            ii.Children[child] = NewIncludeInstructionsEmpty();
            ii.Children[child].Push(inst_rels);
        }
    }
}

func(ii *IncludeInstructions) GetChild(childname string) *IncludeInstructions {
    c := ii.Children[childname];
    if c == nil {
        c = NewIncludeInstructionsEmpty();
    }
    return c;
}
