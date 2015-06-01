package jsonapi;

import (
    . "..";
    . "../extra"
    "net/http"
    "testing"
    "io/ioutil"
);

type Simple struct {
    Id string `jsonapi:"id"`
    Target string
}

func init() {
    sf := NewSessionFactorySimple();

    api := NewAPI(sf);
    
    no_auth := NewAuthenticatorNone();

    resource_simple := NewResourceRAM(&Simple{});
    resource_simple.Push("1", &Simple{Target:"Potato"});
    resource_simple.Push("Potato", &Simple{Target:"1"});

    api.MountResource("source", resource_simple, no_auth);

    api.MountRelationship("rel", "source", "source", NewRelationshipBehaviorFromFieldToId("Target", Required), no_auth);
    go func() {
        err := http.ListenAndServe(":3030", api);
        if(err != nil) {
            panic(err);
        }
    }()
}

func RunTests(t *testing.T, tests map[string]string) {
    for uri,output := range tests {
        resp, err := http.Get("http://localhost:3030"+uri);
        if err != nil {
            t.Fatalf("HTTP Error: %#v\n", err);
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body);
        if(string(body) != output) {
            t.Logf("\nINTENDED OUTPUT: %s\nACTUAL OUTPUT:   %s\nURI: %s\n", output, body, uri);
            t.Fatal();
        }
    }
}

func TestRequestBasics(t *testing.T) {
    tests := map[string]string{
        "/source/1": `{"data":{"attributes":{"Target":"Potato"},"id":"1","relationships":{"rel":{"data":{"type":"source","id":"Potato"},"links":{"related":"http://localhost:3030/source/1/rel","self":"http://localhost:3030/source/1/relationships/rel"}}},"type":"source"}}`,

        "/source/1,Nonexistant": `{"data":[{"attributes":{"Target":"Potato"},"id":"1","relationships":{"rel":{"data":{"type":"source","id":"Potato"},"links":{"related":"http://localhost:3030/source/1/rel","self":"http://localhost:3030/source/1/relationships/rel"}}},"type":"source"}]}`,

        "/source/1,Potato": `{"data":[{"attributes":{"Target":"Potato"},"id":"1","relationships":{"rel":{"data":{"type":"source","id":"Potato"},"links":{"related":"http://localhost:3030/source/1/rel","self":"http://localhost:3030/source/1/relationships/rel"}}},"type":"source"},{"attributes":{"Target":"1"},"id":"Potato","relationships":{"rel":{"data":{"type":"source","id":"1"},"links":{"related":"http://localhost:3030/source/Potato/rel","self":"http://localhost:3030/source/Potato/relationships/rel"}}},"type":"source"}]}`,

        "/source/1/rel": `{"data":{"attributes":{"Target":"1"},"id":"Potato","relationships":{"rel":{"data":{"type":"source","id":"1"},"links":{"related":"http://localhost:3030/source/Potato/rel","self":"http://localhost:3030/source/Potato/relationships/rel"}}},"type":"source"}}`,

        "/source/1?include=rel":`{"data":{"attributes":{"Target":"Potato"},"id":"1","relationships":{"rel":{"data":{"type":"source","id":"Potato"},"links":{"related":"http://localhost:3030/source/1/rel","self":"http://localhost:3030/source/1/relationships/rel"}}},"type":"source"},"included":[{"attributes":{"Target":"1"},"id":"Potato","relationships":{"rel":{"data":{"type":"source","id":"1"},"links":{"related":"http://localhost:3030/source/Potato/rel","self":"http://localhost:3030/source/Potato/relationships/rel"}}},"type":"source"}]}`,

        "/source/1/relationships/rel": `{"data":{"type":"source","id":"Potato"}}`,
    }
    RunTests(t,tests);
}