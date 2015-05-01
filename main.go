package main;

import (
    "./jsonapi"
    "net/http"
    "fmt"
);

type Session struct{
    Id string
}

func(s *Session) GetId() string {
    return s.Id;
}

func(s *Session) SetId(id string) error {
    s.Id = id;
    return nil;
}

type SessionResource struct{
}

func NewSessionResource() *SessionResource {
    return &SessionResource{}
}

func(sr *SessionResource) FindOne(id string, r *http.Request) (jsonapi.HasId, error) {
    return &Session{Id:"123"}, nil;
}

func main() {
    api := jsonapi.NewAPI();

    //api.MountResource("user", jsonapi.NewSQLResource(db, "user"), jsonapi.NoPermissions());
    //api.MountResource("dogs", jsonapi.NewSQLResource(db,"dogs"), jsonapi.NoPermissions());
    api.MountResource("session", NewSessionResource(), jsonapi.NewNoRestrictions());

    //api.MountLinkage("pets", "user", "dogs", jsonapi.SQLLinkageBehavior);

    //api.MountLinkage("login_as", "session", "user", jsonapi.StandardLinkageBehavior);

    // curl localhost:3030/api/user/0/pets

    fmt.Printf("Listening\n");
    err := http.ListenAndServe(":3030", api);
    if err != nil {
        panic(err);
    }
}
