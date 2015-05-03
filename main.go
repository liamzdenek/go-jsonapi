package main;

import (
    "./jsonapi"
    "net/http"
    "fmt"
    //"strconv"
    "time"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
);
type Session struct{
    ID string
    UserId int `json:"user_id"`
    Created time.Time `json:"created,omitempty"`
}

func(s *Session) Id() string {
    return s.ID;
}

/*func(s *Session) SetId(id string) error {
    s.ID = id;
    return nil;
}*/
type SessionResource struct{
}

func NewSessionResource() *SessionResource {
    return &SessionResource{}
}

func(sr *SessionResource) FindMany(ids []string, r *http.Request) ([]jsonapi.Ider, error) {
    return nil, nil;
}

func(sr *SessionResource) FindOne(id string, r *http.Request) (jsonapi.Ider, error) {
    return &Session{ID:"123",Created:time.Now(),UserId:1}, nil;
}
/*
type User struct{
    Id int `meddler:"id,pk"`
    Name string `meddler:"name" json:"name"`
}

func(u *User) GetId() string {
    return fmt.Sprintf("%d",u.Id);
}

func(u *User) SetId(id string) error {
    nid, err := strconv.Atoi(id);
    u.Id = nid;
    return err;
}
*/

func main() {
    _, err := sql.Open("mysql", "root@/tasky");
    if err != nil {
        panic(err);
    }

    api := jsonapi.NewAPI();

    //api.MountResource("user", jsonapi.NewSQLResource(db, "users", &User{}), jsonapi.NewNoRestrictions());
    //api.MountResource("dogs", jsonapi.NewSQLResource(db,"dogs"), jsonapi.NoPermissions());
    api.MountResource("session", NewSessionResource(), jsonapi.NewAuthenticatorNone());

    //api.MountLinkage("pets", "user", "dogs", jsonapi.SQLLinkageBehavior);

    //api.MountLinkage("logged_in_as", "session", "user", jsonapi.NewOneToOneLinkageBehavior("UserId"));

    // curl localhost:3030/api/user/0/pets

    fmt.Printf("Listening\n");
    err = http.ListenAndServe(":3030", api);
    if err != nil {
        panic(err);
    }
}
