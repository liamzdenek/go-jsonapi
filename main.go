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

func(sr *SessionResource) FindMany(ids []string) ([]jsonapi.Ider, error) {
    return []jsonapi.Ider{
        &Session{ID:"123",Created:time.Now(),UserId:1},
        &Session{ID:"124",Created:time.Now(),UserId:2},
        &Session{ID:"125",Created:time.Now(),UserId:3},
    }, nil;
    return nil, nil;
}

func(sr *SessionResource) FindOne(id string) (jsonapi.Ider, error) {
    return &Session{ID:"123",Created:time.Now(),UserId:1}, nil;
}

func(sr *SessionResource) FindManyByField(field string, value string) ([]jsonapi.Ider, error) {
    panic("Session does not support FindManyByField");
}

type User struct{
    ID int `meddler:"id,pk"`
    Name string `meddler:"name" json:"name"`
}

func(u *User) Id() string {
    return fmt.Sprintf("%d",u.ID);
}

type Post struct {
    ID int `meddler:"id,pk"`
    UserId int `meddler:"user_id" json:"-"`
}

func(p *Post) Id() string {
    return fmt.Sprintf("%d",p.ID);
}

type Comment struct {
    ID int `meddler:"id,pk"`
    UserId int `meddler:"user_id" json:"-"`
    PostId int `meddler:"post_id" json:"-"`
    Text string `meddler:"text"`
}

func(c *Comment) Id() string {
    return fmt.Sprintf("%d",c.ID);
}

func main() {
    db, err := sql.Open("mysql", "root@/tasky");
    if err != nil {
        panic(err);
    }

    api := jsonapi.NewAPI();

    api.MountResource("user", jsonapi.NewResourceSQL(db, "users", &User{}), jsonapi.NewAuthenticatorNone());
    api.MountResource("post", jsonapi.NewResourceSQL(db, "posts", &Post{}), jsonapi.NewAuthenticatorNone());
    api.MountResource("comment", jsonapi.NewResourceSQL(db, "comments", &Comment{}), jsonapi.NewAuthenticatorNone());
    api.MountResource("session", NewSessionResource(), jsonapi.NewAuthenticatorNone());

    api.MountRelationship("logged_in_as", "session", "user", jsonapi.NewRelationshipBehaviorFromFieldToId("UserId"), jsonapi.NewAuthenticatorNone());
    api.MountRelationship("posts", "user", "post", jsonapi.NewRelationshipBehaviorFromFieldToField("ID", "UserId"), jsonapi.NewAuthenticatorNone());
    api.MountRelationship("comments", "post", "comment", jsonapi.NewRelationshipBehaviorFromFieldToField("ID", "PostId"), jsonapi.NewAuthenticatorNone());

    // curl localhost:3030/api/user/0/pets
    fmt.Printf("Listening\n");
    err = http.ListenAndServe(":3030", api);
    if err != nil {
        panic(err);
    }
}
