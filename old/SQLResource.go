package jsonapi;

import (
    "database/sql";
    "net/http";
    "github.com/russross/meddler";
    "reflect"
    "strings"
    "fmt"
);

type SQLResource struct{
    DB *sql.DB
    Table string
    Type reflect.Type
}

func init() {
    // safety check to make sure SQLResource is a Resource
    var t Resource;
    t = &SQLResource{};
    _ = t;
}

func NewSQLResource(db *sql.DB, table string, t HasId) *SQLResource {
    return &SQLResource{
        DB: db,
        Table: table,
        Type: reflect.Indirect(reflect.ValueOf(t)).Type(),
    }
}

func(sr *SQLResource) FindOne(id string, r *http.Request) (HasId, error) {
    v := reflect.New(sr.Type).Interface();
    err := meddler.QueryRow(sr.DB, v, "SELECT * FROM "+sr.Table+" WHERE id=?", id);
    return v.(HasId), err;
}

func(sr *SQLResource) FindMany(ids []string, r *http.Request) ([]HasId, error) {
    args := []interface{}{};
    for _, id := range ids {
        args = append(args, id);
    }
    vs := reflect.New(reflect.SliceOf(reflect.PtrTo(sr.Type))).Interface()
    q := "SELECT * FROM "+sr.Table+" WHERE id IN(?"+strings.Repeat(",?", len(ids)-1)+")";
    fmt.Printf("Query: %#v\n", q);
    fmt.Printf("Args: %#v\n", args);
    err := meddler.QueryAll(
        sr.DB,
        vs,
        q,
        args...,
    );
    if(err != nil) {
        return nil, err;
    }
    res := []HasId{};
    fmt.Printf("GOT DATA: %#v\n", vs);

    ary := reflect.Indirect(reflect.ValueOf(vs));
    for i := 0; i < ary.Len(); i++ {
        res = append(res,ary.Index(i).Interface().(HasId));
    }
    fmt.Printf("GOT HASIDS: %#v\n",res);
    return res, err
}
