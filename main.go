package main;

import(. "./jsonapi";"encoding/json";"fmt");

type Test struct {
    ArbitraryField string `json:"qwerty"`
}

func (t *Test) Id() string { return "123"; }

func (t *Test) Link() *OutputLinkageSet { return nil; };

func (t *Test) Type() string { return "test"; };

func main() {
    output := Output{
        Data: &OutputData{
            Data: []*OutputDatum{
                &OutputDatum{
                    Datum: &Test{
                        ArbitraryField: "azerty",
                    },
                },
            },
        },
    };
    bytes, err := json.Marshal(output);
    if err != nil {
        panic(err);
    }
    fmt.Printf("B: %s\n",bytes);
}
