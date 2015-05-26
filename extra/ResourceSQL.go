package jsonapie;

import (
    "database/sql";
    "github.com/russross/meddler";
    "reflect"
    "strings"
    "fmt"
    "errors"
    . ".."
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

func(sr *ResourceSQL) CastSession(s Session) SessionResourceSQL {
    // TODO: proper error handling
    return s.(SessionResourceSQL);
}

// TODO: update this to honor sorting
func(sr *ResourceSQL) FindDefault(a *API, s Session, p *Paginator) ([]Ider, error) {
    vs := reflect.New(reflect.SliceOf(reflect.PtrTo(sr.Type))).Interface()
    offset_and_limit := "";
    if p != nil && (p.MaxPerPage != 0) {
        offset_and_limit = fmt.Sprintf("LIMIT %d OFFSET %d",
            p.MaxPerPage,
            p.CurPage * p.MaxPerPage,
        );
    }
    q := fmt.Sprintf(
        "SELECT * FROM %s %s",
        sr.Table,
        offset_and_limit,
    )
    a.Logger.Printf("Query: %#v\n", q);
    err := meddler.QueryAll(
        sr.DB,
        vs,
        q,
    );
    if(err != nil) {
        return nil, err;
    }
    return sr.ConvertInterfaceSliceToIderSlice(vs), err
;
}

func(sr *ResourceSQL) FindOne(a *API, s Session, id string) (Ider, error) {
    v := reflect.New(sr.Type).Interface();
    err := meddler.QueryRow(sr.DB, v, "SELECT * FROM "+sr.Table+" WHERE id=?", id);
    return v.(Ider), err;
}

func(sr *ResourceSQL) FindMany(a *API, s Session, p *Paginator, ids []string) ([]Ider, error) {
    args := []interface{}{};
    for _, id := range ids {
        args = append(args, id);
    }
    vs := reflect.New(reflect.SliceOf(reflect.PtrTo(sr.Type))).Interface()
    offset_and_limit := "";
    if p != nil {
        offset_and_limit = fmt.Sprintf("LIMIT %d OFFSET %d",
            (*p).MaxPerPage,
            (*p).CurPage * (*p).MaxPerPage,
        );
    }
    q := fmt.Sprintf(
        "SELECT * FROM %s WHERE id IN(?"+strings.Repeat(",?", len(ids)-1)+") %s",
        sr.Table,
        offset_and_limit,
    );

    a.Logger.Printf("Query: %#v\n", q);
    //a.Logger.Printf("Args: %#v\n", args);
    err := meddler.QueryAll(
        sr.DB,
        vs,
        q,
        args...,
    );
    if(err != nil) {
        return nil, err;
    }
    return sr.ConvertInterfaceSliceToIderSlice(vs), err
}

func(sr *ResourceSQL) FindManyByField(a *API, s Session, field string, value string) ([]Ider, error) {
    vs := reflect.New(reflect.SliceOf(reflect.PtrTo(sr.Type))).Interface();
    field, err := sr.GetTableFieldFromStructField(field);
    if(err != nil) {
        return nil, err;
    }
    // TODO: find a way to parameterize field in this query
    // right now, field is always a trusted string, but some
    // later relationship behaviors might change that, and it's
    // better to be safe than sorry
    // dropping in ? instead of field does not work :/
    q := "SELECT * FROM "+sr.Table+" WHERE "+field+"=?";
    a.Logger.Printf("Query: %#v %#v\n", q, value);
    err = meddler.QueryAll(
        sr.DB,
        vs,
        q,
        value,
    );
    //a.Logger.Printf("RES: %#v\n", vs);
    return sr.ConvertInterfaceSliceToIderSlice(vs), err;
}

func(sr *ResourceSQL) Delete(a *API, s Session, id string) error {
    _, err := sr.DB.Exec("DELETE FROM "+sr.Table+" WHERE id=?", id);
    return err;
}


func(sr *ResourceSQL) ParseJSON(a *API, s Session, ider Ider, raw []byte) (Ider, *string, *string, *OutputLinkageSet, error) {
    return ParseJSONHelper(ider, raw, sr.Type);
}

func(sr *ResourceSQL) Create(a *API, s Session, ider Ider, id *string) (RecordCreatedStatus, error) {
    sqlctx := sr.CastSession(s);
    a.Logger.Printf("CREATE GOT CONTEXT: %#v\n", sqlctx);
    if(id != nil) {
        return StatusFailed, errors.New("ResourceSQL does not support specifying an ID for Create() requests."); // TODO: it should
    }
    tx, err := sqlctx.GetSQLTransaction(sr.DB)
    if err != nil {
        return StatusFailed, err;
    }
    err = meddler.Insert(tx, sr.Table, ider)
    return StatusCreated, err;
}

func(sr *ResourceSQL) Update(a *API, s Session, id string, ider Ider) error {
    panic("NOT IMPLEMENTED");
}

func (sr *ResourceSQL) ConvertInterfaceSliceToIderSlice(src interface{}) []Ider {
    res := []Ider{};

    ary := reflect.Indirect(reflect.ValueOf(src));
    for i := 0; i < ary.Len(); i++ {
        res = append(res,ary.Index(i).Interface().(Ider));
    }
    return res;
}

func (sr *ResourceSQL) GetTableFieldFromStructField(structstr string) (string, error) {
    field, found := sr.Type.FieldByName(structstr);
    if(!found) {
        return "", errors.New("Field "+structstr+" does not exist on "+sr.Type.Name());
    }
    realname := field.Name;

    meddler_tags := strings.Split(field.Tag.Get("meddler"),",");

    if(len(meddler_tags) > 0 && meddler_tags[0] != "") {
        realname = meddler_tags[0];
    }

    return realname, nil;
}
