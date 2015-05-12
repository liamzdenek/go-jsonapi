package main;

import (
    . ".."
    . "../extra"
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

func(sr *SessionResource) FindMany(ids []string) ([]Ider, error) {
    return []Ider{
        &Session{ID:"123",Created:time.Now(),UserId:1},
        &Session{ID:"124",Created:time.Now(),UserId:2},
        &Session{ID:"125",Created:time.Now(),UserId:3},
    }, nil;
    return nil, nil;
}

func(sr *SessionResource) FindOne(id string) (Ider, error) {
    return &Session{ID:"123",Created:time.Now(),UserId:1}, nil;
}

func(sr *SessionResource) FindManyByField(field string, value string) ([]Ider, error) {
    panic(NewResponderErrorOperationNotSupported("Session does not support FindManyByField"));
}

func(sr *SessionResource) Delete(id string) error {
    panic(NewResponderErrorOperationNotSupported("Session does not support Delete"));
}

func(sr *SessionResource) Create(resource_str string, raw []byte, verify ResourceCreateVerifyFunc) (Ider, RecordCreatedStatus, error) {
    panic(NewResponderErrorOperationNotSupported("Session does not support Create"));
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
    Text string `meddler:"text" json:"text"`
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

    api := NewAPI();

    api.MountResource("user", NewResourceSQL(db, "users", &User{}), NewAuthenticatorNone());
    api.MountResource("post", NewResourceSQL(db, "posts", &Post{}), NewAuthenticatorNone());
    api.MountResource("comment", NewResourceSQL(db, "comments", &Comment{}), NewAuthenticatorNone());
    api.MountResource("session", NewSessionResource(), NewAuthenticatorNone());

    api.MountRelationship("logged_in_as", "session", "user",
        NewRelationshipBehaviorFromFieldToId("UserId", Required),
        NewAuthenticatorNone(),
    );
    api.MountRelationship("posts", "user", "post",
        NewRelationshipBehaviorFromFieldToField("ID", "UserId", Required),
        NewAuthenticatorNone(),
    );
    api.MountRelationship("author", "post", "user",
        NewRelationshipBehaviorFromFieldToId("UserId", Required),
        NewAuthenticatorNone(),
    );
    api.MountRelationship("comments", "post", "comment",
        NewRelationshipBehaviorFromFieldToField("ID", "PostId", Required),
        NewAuthenticatorNone(),
    );

    // curl localhost:3030/api/user/0/pets
    fmt.Printf("Listening\n");
    err = http.ListenAndServe(":3030", api);
    if err != nil {
        panic(err);
    }
}
