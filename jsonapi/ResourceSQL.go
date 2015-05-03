package jsonapi;

import (
    "database/sql";
    "github.com/russross/meddler";
    "reflect"
    "strings"
    "fmt"
);

type ResourceSQL struct{
    DB *sql.DB
    Table string
    Type reflect.Type
}

func init() {
    // safety check to make sure ResourceSQL is a Resource
    var t Resource;
    t = &ResourceSQL{};
    _ = t;
}

func NewResourceSQL(db *sql.DB, table string, t Ider) *ResourceSQL {
    return &ResourceSQL{
        DB: db,
        Table: table,
        Type: reflect.Indirect(reflect.ValueOf(t)).Type(),
    }
}

func(sr *ResourceSQL) FindOne(id string) (Ider, error) {
    v := reflect.New(sr.Type).Interface();
    err := meddler.QueryRow(sr.DB, v, "SELECT * FROM "+sr.Table+" WHERE id=?", id);
    return v.(Ider), err;
}

func(sr *ResourceSQL) FindMany(ids []string) ([]Ider, error) {
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
    res := []Ider{};
    fmt.Printf("GOT DATA: %#v\n", vs);

    ary := reflect.Indirect(reflect.ValueOf(vs));
    for i := 0; i < ary.Len(); i++ {
        res = append(res,ary.Index(i).Interface().(Ider));
    }
    fmt.Printf("GOT HASIDS: %#v\n",res);
    return res, err
}
