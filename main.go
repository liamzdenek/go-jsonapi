package main;

import (
    "./jsonapi"
    "net/http"
);

func main() {
    api := jsonapi.NewAPI();

    api.MountResource("user", jsonapi.NewSQLResource(db, "user"), jsonapi.NoPermissions());
    api.MountResource("dogs", jsonapi.NewSQLResource(db,"dogs"), jsonapi.NoPermissions());
    api.MountResource("session", NewSessionResource(), jsonapi.NoPermissions());

    api.MountLinkage("pets", "user", "dogs", jsonapi.SQLLinkageBehavior);

    api.MountLinkage("login_as", "session", "user", jsonapi.StandardLinkageBehavior);

    // curl localhost:3030/api/user/0/pets

    err := http.ListenAndServe(":3030", api);
    if err != nil {
        panic(err);
    }
}
