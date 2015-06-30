package main;

import (
    . ".."
    "../resource"
    "../relationship"
    //. "../oauth2"
    //"../authenticator"
    "net/http"
    "fmt"
    //"strconv"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "runtime"
    //"time"
);

/*
func init() {
    go func() {
        counter := 0;
        for {
            time.Sleep(time.Second);
            fmt.Printf("GOROUTINES: %#v\n", runtime.NumGoroutine());
            counter++;
            if counter > 10 {
                panic("10 seconds has passed");
            }
        }
    }();
}
*/

type User struct{
    ID int `meddler:"id,pk" jsonapi:"id" json:"-"`
    Name string `meddler:"name" json:"name"`
}

type Post struct {
    ID int `meddler:"id,pk" jsonapi:"id"`
    UserId int `meddler:"user_id" json:"-"`
    Text string `meddler:"text" json:"text"`
}

type Comment struct {
    ID int `meddler:"id,pk" jsonapi:"id"`
    UserId int `meddler:"user_id" json:"-"`
    PostId int `meddler:"post_id" json:"-"`
    Text string `meddler:"text"`
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU() * 4)
    //base_oauth := "/auth/";
    base_api := "/api/";

    db, err := sql.Open("mysql", "root@/tasky");
    if err != nil {
        panic(err);
    }

    api := NewAPI(base_api);

    //oauth2 := NewOAuth2(base_oauth);

    resource_user := resource.NewSQL(db, "users", &User{}) // database, table name, raw struct to unwrap into
    resource_post := resource.NewSQL(db, "posts", &Post{})
    resource_comment := resource.NewSQL(db, "comments", &Comment{})
    /*
    resource_session := resource.NewSQL(db, "session", &authenticator.SimpleLoginSession{});
    
    resource_login := authenticator.NewSimpleLogin(resource_user,resource_session);
    no_auth := authenticator.NewNone();

    rbac := authenticator.NewRBAC(
        resource.NewSQL(db, "rbac_permissions", &authenticator.RBACPermissionLookup{}),
        resource.NewSQL(db, "rbac_user_permissions", &authenticator.RBACUserPermissionLookup{}),
        resource_login,
    );
    */

    api.MountResource("user", resource_user);
    api.MountResource("post", resource_post);
    api.MountResource("comment", resource_comment/*, rbac.Require("canComment")*/);
    /*
    api.MountResource("login", resource_login);

    */
    api.MountRelationship("posts", "user", "post",
        relationship.NewFromFieldToField("ID", "UserId"),
    );
    api.MountRelationship("author", "post", "user",
        relationship.NewFromFieldToId("UserId"),
    );
    api.MountRelationship("comments", "post", "comment",
        relationship.NewFromFieldToField("ID", "PostId"),
    );

    // that's it! start the API
    fmt.Printf("Listening\n");
    //http.Handle(base_oauth, oauth2);
    http.Handle(base_api, api);
    err = http.ListenAndServe(":3030", nil);
    if err != nil {
        panic(err);
    }
}
