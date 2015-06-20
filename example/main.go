package main;

import (
    . ".."
    "../resource"
    "../relationship"
    //. "../oauth2"
    "../authenticator"
    "net/http"
    "fmt"
    //"strconv"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
);

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
    resource_session := resource.NewSQL(db, "session", &authenticator.SimpleLoginSession{});

    resource_login := authenticator.NewSimpleLogin(resource_user,resource_session);
    no_auth := authenticator.NewNone();

    rbac := authenticator.NewRBAC(
        resource.NewSQL(db, "rbac_permissions", &authenticator.RBACPermissionLookup{}),
        resource.NewSQL(db, "rbac_user_permissions", &authenticator.RBACUserPermissionLookup{}),
        resource_login,
    );

    api.MountResource("user", resource_user);
    api.MountResource("post", resource_post);
    api.MountResource("comment", resource_comment, rbac.Require("canComment"));
    api.MountResource("login", resource_login);

    // Relationships are the other large concept in addition to Resources. Each relationship
    // is a unidirectional associations between two resources. The relationship is given a name--
    // in the following example, the relationship name is "logged_in_as", and it associates a given
    // "session" with a "user." 
    //
    // Relationships have a Relationship Behavior, which is responsible for taking any arbitrary
    // record of the source resource, and converting it into the destination resource. There are
    // no behaviors built in to core, but there are a few in jsonapi/extras
    //
    // Finally, relationships have an authenticator, which can refuse access upon any arbitrary
    // conditions, just like Resources TODO: this isn't totally working as desired yet, so use with caution
    // Example requests (note, these are not exact responses as they have been simplified to exaggerate mechanics):
    // * curl -X GET "localhost:3030/session/1/logged_in_as"
    //   * returns: {"data":{"type":"user","id":"1","attributes":{"name":"Jsonapi IsGreat"}}}
    //
    // * curl -X GET "localhost:3030/session/1/links/logged_in_as
    //   * returns: {"data":{"type":"user","id":"1"}}
    /*
    api.MountRelationship("logged_in_as", "session", // name, src, dest
        NewRelationshipFromFieldToId("user", "UserId", Required),
        no_auth,
    );
    */
    api.MountRelationship("posts", "user",
        relationship.NewFromFieldToField("post", "ID", "UserId", NotRequired),
        no_auth,
    );
    api.MountRelationship("author", "post",
        relationship.NewFromFieldToId("user", "UserId", Required),
        no_auth,
    );
    api.MountRelationship("comments", "post",
        relationship.NewFromFieldToField("comment", "ID", "PostId", Required),
        no_auth,
    );

    // For requests with a resource as the primary "data" (aka, not a /:resource/:id/links/:linkname request),
    // you may specify ?include=[includefmt] to include the data for additional links. For example:
    //
    // * curl -X GET "localhost:3030/session/1/logged_in_as?include=posts,posts.comments
    //   * includes the additional field "includes" populated with all the data needed to navigate through the links.

    // For guided documentation, see the 'docs' folder

    // that's it! start the API
    fmt.Printf("Listening\n");
    //http.Handle(base_oauth, oauth2);
    http.Handle(base_api, api);
    err = http.ListenAndServe(":3030", nil);
    if err != nil {
        panic(err);
    }
}
