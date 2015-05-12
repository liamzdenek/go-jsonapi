package main;

import(. "./jsonapi";"encoding/json";"fmt");

type Session struct {
    ArbitraryField string `json:"qwerty"`
}

func (s *Session) Id() string { return "123"; }


func (s *Session) Type() string { return "session"; }
func (s *Session) Link() *OutputLinkageSet { return &OutputLinkageSet{
    Linkages: []*OutputLinkage{
        &OutputLinkage{
            LinkName: "logged_in_as",
            Links: []OutputLink{
                {
                    Type:"user",
                    Id: "123",
                },
            },
        },
    },
}};

func main() {
    t :=  &Session{
        ArbitraryField: "azerty",
    };
    output := Output{}
    output.Data = NewOutputDataResources(false, []*OutputDatum{
        &OutputDatum{
            Datum: t,
        },
    });
    bytes, _ := json.Marshal(output);
    fmt.Printf("A: %s\n",bytes);

    output.Data = NewOutputDataRelationship(t.Link());
    bytes, _ = json.Marshal(output);
    fmt.Printf("B: %s\n",bytes);

    output.Data = NewOutputDataLinkage(true, t.Link().Linkages[0]);
    bytes, _ = json.Marshal(output);
    fmt.Printf("C: %s\n",bytes);
}
