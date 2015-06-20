package main;

import (
    . ".."
    . "../extra"
    . "../oauth2"
    "../authenticator"
    "net/http"
    "fmt"
    //"strconv"
    "time"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
);
type UserSession struct{
    ID string `jsonapi:"id"`
    UserId int `json:"-"`
    Created *time.Time `json:"created"`
}

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

type SimpleGetLoggedInAs struct {}

func(s *SimpleGetLoggedInAs) GetUserId(r *Request) *string {
    user_id := "9";
    return &user_id;
}

func main() {
    base_oauth := "/auth/";
    base_api := "/api/";

    db, err := sql.Open("mysql", "root@/tasky");
    if err != nil {
        panic(err);
    }

    api := NewAPI(base_api);

    oauth2 := NewOAuth2(base_oauth);

    // Resources are one of the primary concepts in jsonapi. A resource defines CRUD primitives for
    // retrieving, manipulating, and parsing the underlying data. There are no resources built in to core,
    // but the 'extras' folder (jsonapie) provides a few, such as ResourceSQL
    resource_user := NewResourceSQL(db, "users", &User{}) // database, table name, raw struct to unwrap into
    resource_post := NewResourceSQL(db, "posts", &Post{})
    resource_comment := NewResourceSQL(db, "comments", &Comment{})
    //resource_session := NewResourceRAM(&UserSession{});

    // load up some test data to the session since it is entirely in RAM
    //now := time.Now();
    //resource_session.Push("1", &UserSession{ID: "1", UserId: 1, Created:&now});
    //resource_session.Push("2", &UserSession{ID: "2", UserId: 2, Created:&now});
    //resource_session.Push("3", &UserSession{ID: "3", UserId: 17, Created:&now});

    // initialize our authentication scheme. The authenticator is where you put code to
    // permit or refuse access to certain resources or linkages based on whatever rules
    // for example, a resource might be admin-only, or retrieving a certain link might
    // depend on any arbitrary condition
    // AuthenticatorNone() is the only built in authenticator. Every request is granted.
    no_auth := NewAuthenticatorNone();
    rbac := authenticator.NewRBAC(
        NewResourceSQL(db, "rbac_permissions", &authenticator.RBACPermissionLookup{}),
        NewResourceSQL(db, "rbac_user_permissions", &authenticator.RBACUserPermissionLookup{}),
        &SimpleGetLoggedInAs{},
    );

    // Resources can be easily wrapped with common functionality,
    // such as caching, or pagination. These are designed
    // to chain easily (not shown below)
    //resource_comment_cache := NewResourceCache("posts", resource_comment, NewResourceRAM(&CacheRecord{}))
    //resource_post_paginator := NewResourcePaginatorSimple(5, resource_post);
    resource_comment_cache := resource_comment;
    resource_post_paginator := resource_post;

    // api.MountResource informs the api of the provided resource, and makes the resource
    // available to requests using the string given as the first parameter.
    // Example requests (note, these are not exact responses as they have been simplified to exaggerate mechanics):
    // * curl -X GET "localhost:3030/user/1"
    //   * returns: {"data":{"type":"user","id":"1","attributes":{"name":"Jsonapi IsGreat"}}}
    //
    // * curl -X GET "localhost:3030/user/1,2,3"
    //   * returns: {"data":[{"type":"user","id":"1","attributes":{"name":"Jsonapi IsGreat"}},{"type":"user","id":2", ...} ...]}
    //
    // * curl -X POST "localhost:3030/user" -d '{"data":"{"type":"user","id":"2","attributes":{"name":"JsonapiMakes Apis SoEasy"}}}'
    //   * returns: 201 Created + Location header upon success
    // TODO: write real docs about how all of this works
    api.MountResource("user", resource_user, no_auth);
    api.MountResource("post", resource_post_paginator, no_auth);
    api.MountResource("comment", resource_comment_cache, rbac.Require("canComment"));
    //api.MountResource("session", resource_session, no_auth);

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
        NewRelationshipFromFieldToField("post", "ID", "UserId", NotRequired),
        no_auth,
    );
    api.MountRelationship("author", "post",
        NewRelationshipFromFieldToId("user", "UserId", Required),
        no_auth,
    );
    api.MountRelationship("comments", "post",
        NewRelationshipFromFieldToField("comment", "ID", "PostId", Required),
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
    http.Handle(base_oauth, oauth2);
    http.Handle(base_api, api);
    err = http.ListenAndServe(":3030", nil);
    if err != nil {
        panic(err);
    }
}
