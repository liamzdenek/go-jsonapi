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

func init() {
    // safety check to make sure ResourceSQL is a Resource
    var t Resource;
    t = &ResourceSQL{};
    _ = t;
}

type ResourceSQL struct{
    DB *sql.DB
    Table string
    Type reflect.Type
}

type ResourceSQLPromise struct {
    Transactions map[*sql.DB]*sql.Tx;
}

func(rsp *ResourceSQLPromise) GetSQLTransaction(db *sql.DB) (*sql.Tx, error) {
    if tx, ok := rsp.Transactions[db]; ok && tx != nil {
        return tx, nil;
    }
    res,err := db.Begin();
    if err != nil {
        return nil, err;
    }
    rsp.Transactions[db] = res;
    return res, nil;
}

func(rsp *ResourceSQLPromise) Failure(r *Request) {
    r.API.Logger.Infof("ResourceSQLPromise Failure\n");
    for _,tx := range rsp.Transactions {
        err := tx.Rollback();
        Check(err);
    }
}

func(rsp *ResourceSQLPromise) Success(r *Request) {
    r.API.Logger.Infof("ResourceSQLPromise Success\n");
    for _,tx := range rsp.Transactions {
        err := tx.Commit();
        Check(err);
    }
}

func NewResourceSQL(db *sql.DB, table string, t interface{}) *ResourceSQL {
    return &ResourceSQL{
        DB: db,
        Table: table,
        Type: reflect.Indirect(reflect.ValueOf(t)).Type(),
    }
}

func (sr *ResourceSQL) GetPromise(r *Request) (LeasedPromise, *ResourceSQLPromise) {
    v := r.PromiseStorage.Get(&ResourceSQLPromise{}, func() Promise {
        return &ResourceSQLPromise{
            Transactions: make(map[*sql.DB]*sql.Tx),
        };
    });
    return v, v.Promise.(*ResourceSQLPromise);
}

// TODO: update this to honor sorting
func(sr *ResourceSQL) FindDefault(r *Request, rp RequestParams) ([]*Record, error) {
    p := rp.Paginator;
    a := r.API;
    vs := reflect.New(reflect.SliceOf(reflect.PtrTo(sr.Type))).Interface()
    offset_and_limit := "";
    if p.ShouldPaginate {
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
    a.Logger.Debugf("Query: %#v\n", q);
    err := meddler.QueryAll(
        sr.DB,
        vs,
        q,
    );
    if(err != nil) {
        return nil, err;
    }
    return sr.ConvertInterfaceSliceToRecordSlice(vs), err;
}

func(sr *ResourceSQL) FindOne(r *Request, rp RequestParams, id string) (*Record, error) {
    v := reflect.New(sr.Type).Interface();
    err := meddler.QueryRow(sr.DB, v, "SELECT * FROM "+sr.Table+" WHERE id=?", id);
    if err != nil {
        return nil, err;
    }
    return &Record{
        Id: id,
        Attributes: v,
    }, nil;
}
func(sr *ResourceSQL) FindMany(r *Request, rp RequestParams, ids []string) ([]*Record, error) {
    p := rp.Paginator;
    args := []interface{}{};
    for _, id := range ids {
        args = append(args, id);
    }
    vs := reflect.New(reflect.SliceOf(reflect.PtrTo(sr.Type))).Interface()
    offset_and_limit := "";
    if p.ShouldPaginate {
        offset_and_limit = fmt.Sprintf("LIMIT %d OFFSET %d",
            p.MaxPerPage,
            p.CurPage * p.MaxPerPage,
        );
    }
    q := fmt.Sprintf(
        "SELECT * FROM %s WHERE id IN(?"+strings.Repeat(",?", len(ids)-1)+") %s",
        sr.Table,
        offset_and_limit,
    );

    r.API.Logger.Debugf("Query: %#v\n", q);
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
    r.API.Logger.Debugf("QUERY RES: %#v\n", vs);
    return sr.ConvertInterfaceSliceToRecordSlice(vs), err
}

func(sr *ResourceSQL) FindManyByField(r *Request, rp RequestParams, field, value string) ([]*Record, error) {
    p := rp.Paginator;
    vs := reflect.New(reflect.SliceOf(reflect.PtrTo(sr.Type))).Interface();
    field, err := sr.GetTableFieldFromStructField(field);
    if(err != nil) {
        return nil, err;
    }
    offset_and_limit := "";
    if p.ShouldPaginate {
        offset_and_limit = fmt.Sprintf("LIMIT %d OFFSET %d",
            p.MaxPerPage,
            p.CurPage * p.MaxPerPage,
        );
    }
    // TODO: find a way to parameterize field in this query
    // right now, field is always a trusted string, but some
    // later relationship behaviors might change that, and it's
    // better to be safe than sorry
    // dropping in ? instead of field does not work :/
    q := fmt.Sprintf("SELECT * FROM %s WHERE %s=? %s", sr.Table, field, offset_and_limit);
    r.API.Logger.Debugf("Query: %#v %#v\n", q, value);
    err = meddler.QueryAll(
        sr.DB,
        vs,
        q,
        value,
    );
    r.API.Logger.Debugf("RES: %#v\n", vs);
    return sr.ConvertInterfaceSliceToRecordSlice(vs), err;
}

func(sr *ResourceSQL) Delete(r *Request, id string) error {
    _, err := sr.DB.Exec("DELETE FROM "+sr.Table+" WHERE id=?", id);
    return err;
}


func(sr *ResourceSQL) ParseJSON(r *Request, src *Record, raw []byte) (*Record, error) {
    return ParseJSONHelper(src, raw, sr.Type);
}

func(sr *ResourceSQL) Create(r *Request, src *Record) (RecordCreatedStatus, error) {
    lp, psql := sr.GetPromise(r);
    defer lp.Release();
    r.API.Logger.Debugf("CREATE GOT CONTEXT: %#v\n", psql);
    if(src.Id != "") {
        return StatusFailed, errors.New("ResourceSQL does not support specifying an ID for Create() requests."); // TODO: it should
    }
    tx, err := psql.GetSQLTransaction(sr.DB)
    if err != nil {
        return StatusFailed, err;
    }
    SetId(src.Attributes, src.Id);
    err = meddler.Insert(tx, sr.Table, src.Attributes)
    if err != nil {
        return StatusFailed, err;
    }
    return StatusCreated, err;
}

func(sr *ResourceSQL) Update(r *Request, rec *Record) error {
    panic("NOT IMPLEMENTED");
}

func (sr *ResourceSQL) ConvertInterfaceSliceToRecordSlice(src interface{}) []*Record {
    res := []*Record{};

    ary := reflect.Indirect(reflect.ValueOf(src));
    for i := 0; i < ary.Len(); i++ {
        attr := ary.Index(i).Interface();
        res = append(res,&Record{
            Id: GetId(attr),
            Attributes: attr,
        });
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
