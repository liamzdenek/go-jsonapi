package resource;

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
    // safety check to make sure SQL is a Resource
    var t Resource;
    t = &SQL{};
    _ = t;
}

type SQL struct{
    DB *sql.DB
    Table string
    Type reflect.Type
}

type SQLPromise struct {
    Transactions map[*sql.DB]*sql.Tx;
}

func(rsp *SQLPromise) GetSQLTransaction(db *sql.DB) (*sql.Tx, error) {
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

func(rsp *SQLPromise) Failure(r *Request) {
    r.API.Logger.Infof("SQLPromise Failure\n");
    for _,tx := range rsp.Transactions {
        err := tx.Rollback();
        Check(err);
    }
}

func(rsp *SQLPromise) Success(r *Request) {
    r.API.Logger.Infof("SQLPromise Success\n");
    for _,tx := range rsp.Transactions {
        err := tx.Commit();
        Check(err);
    }
}


type FutureSQL struct {
    Resource *SQL
    Children []*FutureSQL
    Limit, Offset uint
}

func(f *FutureSQL) PrepareQuery(parameters ...SQLExpression) (query string, arguments []interface{}, is_single bool) {
    q := &SQLQuery{
        Query: "SELECT * FROM %s ",
        FmtArguments: []interface{}{f.Resource.Table},
    };
    if len(parameters) > 0 {
        params := NewSQLWhere(
            NewSQLAnd(
                parameters...,
            ),
        );
        if params != nil {
            params.Express(q);
        }
    }
    query = q.PrepareQuery();
    fmt.Printf("Got query: %s\n", query);
    return query, q.SqlArguments, false; // TODO
}

func(f *FutureSQL) RunQuery(pf *ExecutableFuture, req *FutureRequest, parameters []SQLExpression) ([]*Record, bool, error){
    vs := reflect.New(reflect.SliceOf(reflect.PtrTo(f.Resource.Type))).Interface()
    query, queryargs, is_single := f.PrepareQuery(parameters...);
    pf.Request.API.Logger.Debugf("RUN QUERY: %#v %#v\n", query, queryargs);
    err := meddler.QueryAll(
        f.Resource.DB,
        vs,
        query,
        queryargs...,
    );
    return f.Resource.ConvertInterfaceSliceToRecordSlice(vs), is_single, err;
}

func(f *FutureSQL) WorkFindByFields(pf *ExecutableFuture, req *FutureRequest, k *FutureRequestKindFindByAnyFields) {
    parameters := []SQLExpression{};
    forced_single := false;
    if len(k.Fields) > 0 {
        forced_single = len(k.Fields) == 1;
        for _, field := range k.Fields {
            //field_key := f.Resource.GetIdFieldName(nil);
            parameters = append(parameters, &SQLEquals{
                Field: f.Resource.GetFieldByName(nil, field.Field),
                Value: field.Value,
            });
        }
        parameters = []SQLExpression{
            NewSQLOr(parameters...),
        }
    }
    records, is_single, err := f.RunQuery(pf, req, parameters);
    if err != nil {
        req.SendResponse(&FutureResponse{
            IsSuccess: false,
            Failure: []OError{ErrorToOError(err)},
        });
        return;
    }
    field_records := &FutureResponseKindByFields{
        Records: map[Field][]*Record{},
    };
    for _, record := range records {
        for _, field := range k.Fields {
            // TODO: probably a better way to do this somehow
            if(record.HasFieldValue(field)) {
                field_records.Records[field] = append(field_records.Records[field], record);
            }
        }
    }
    pf.Request.API.Logger.Debugf("GOT RECORDS: %#v\n", records);
    pf.Request.API.Logger.Debugf("SENDING BACK FR: %#v\n", field_records);
    _ = forced_single;
    _ = is_single;
    res := &FutureResponse{
        IsSuccess: true,
        Success: map[Future]FutureResponseKind{
            f: field_records,
        },
    }
    req.SendResponse(res);
}

func(f *FutureSQL) WorkFindByIds(pf *ExecutableFuture, req *FutureRequest, k *FutureRequestKindFindByIds) {
    parameters := []SQLExpression{};
    forced_single := false;
    if len(k.Ids) > 0 {
        forced_single = len(k.Ids) == 1;
        id_field := f.Resource.GetIdFieldName(nil);
        for _, id := range k.Ids {
            parameters = append(parameters, &SQLEquals{
                Field: id_field,
                Value: id,
            });
        }
        parameters = []SQLExpression{
            NewSQLOr(parameters...),
        }
    }
    records, is_single, err := f.RunQuery(pf, req, parameters);
    if err != nil {
        req.SendResponse(&FutureResponse{
            IsSuccess: false,
            Failure: []OError{ErrorToOError(err)},
        });
        return;
    }
    req.SendResponse(&FutureResponse{
        IsSuccess: true,
        Success: map[Future]FutureResponseKind{
            f: &FutureResponseKindRecords{
                IsSingle: forced_single || is_single,
                Records: records,
            },
        },
    });
}

func(f *FutureSQL) Work(pf *ExecutableFuture) {
    for {
        req := pf.GetRequest();
        switch k := req.Kind.(type) {
        case *FutureRequestKindFindByIds:
            f.WorkFindByIds(pf,req,k);
        case *FutureRequestKindFindByAnyFields:
            f.WorkFindByFields(pf,req,k);
        default:
            panic(fmt.Sprintf("FutureSQL got unsupported query kind %T: %#v\n", req.Kind, req.Kind));
        }
    }
}

func(f *FutureSQL) ShouldCombine(n Future) bool { return false; }
func(f *FutureSQL) Combine(n Future) error { return nil;}

func NewSQL(db *sql.DB, table string, t interface{}) *SQL {
    return &SQL{
        DB: db,
        Table: table,
        Type: reflect.Indirect(reflect.ValueOf(t)).Type(),
    }
}

func (sr *SQL) GetPromise(r *Request) (LeasedPromise, *SQLPromise) {
    v := r.PromiseStorage.Get(&SQLPromise{}, func() Promise {
        return &SQLPromise{
            Transactions: make(map[*sql.DB]*sql.Tx),
        };
    });
    return v, v.Promise.(*SQLPromise);
}

func(sr *SQL) GetIdFieldName(v interface{}) string {
    if v == nil {
        v = reflect.New(sr.Type).Interface();
    }
    _, id_field := GetIdField(v);
    id_sql_name := id_field.Name;
    if meddler_tag := id_field.Tag.Get("meddler"); len(meddler_tag) > 0 {
        parts := strings.Split(meddler_tag,",");
        id_sql_name = parts[0];
    }
    return id_sql_name;
}

func(sr *SQL) GetFieldByName(v interface{}, field string) string {
    if v == nil {
        v = reflect.New(sr.Type).Interface();
    }
    val := reflect.Indirect(reflect.ValueOf(v));
    typ := val.Type();
    fields := val.NumField();
    for i := 0; i < fields; i++ {
        if typ.Field(i).Name == field {
            tags := strings.Split(typ.Field(i).Tag.Get("meddler"),",");
            return tags[0];
        }
    }
    return field;
}

func(sr *SQL) GetFuture() Future {
    return &FutureSQL{
        Resource: sr,
    }
}

/*
// TODO: update this to honor sorting
func(sr *SQL) FindDefault(r *Request, rp RequestParams) ([]*Record, error) {
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
    if err == sql.ErrNoRows {
        return nil, nil;
    }
    if(err != nil) {
        return nil, err;
    }
    return sr.ConvertInterfaceSliceToRecordSlice(vs), err;
}

func(sr *SQL) FindOne(r *Request, rp RequestParams, id string) (*Record, error) {
    v := reflect.New(sr.Type).Interface();
    id_sql_name := sr.GetIdFieldName(v);
    err := meddler.QueryRow(sr.DB, v, "SELECT * FROM "+sr.Table+" WHERE "+id_sql_name+"=?", id);
    if err == sql.ErrNoRows {
        return nil, nil;
    }
    if err != nil {
        return nil, err;
    }
    return &Record{
        Id: id,
        Attributes: v,
    }, nil;
}
func(sr *SQL) FindMany(r *Request, rp RequestParams, ids []string) ([]*Record, error) {
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
    id_sql_name := sr.GetIdFieldName(nil);
    q := fmt.Sprintf(
        "SELECT * FROM %s WHERE %s IN(?"+strings.Repeat(",?", len(ids)-1)+") %s",
        sr.Table,
        id_sql_name,
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
    if err == sql.ErrNoRows {
        return nil, nil;
    }
    if(err != nil) {
        return nil, err;
    }
    r.API.Logger.Debugf("QUERY RES: %#v\n", vs);
    return sr.ConvertInterfaceSliceToRecordSlice(vs), err
}

func(sr *SQL) FindManyByField(r *Request, rp RequestParams, field, value string) ([]*Record, error) {
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
    if err == sql.ErrNoRows {
        return nil, nil;
    }
    r.API.Logger.Debugf("RES: %#v\n", vs);
    return sr.ConvertInterfaceSliceToRecordSlice(vs), err;
}
*/

func(sr *SQL) Delete(r *Request, id string) {
    id_sql_name := sr.GetIdFieldName(nil);
    _, err := sr.DB.Exec("DELETE FROM "+sr.Table+" WHERE "+id_sql_name+"=?", id);
    if err != nil {
        panic(err);
    }
}


func(sr *SQL) ParseJSON(r *Request, src *Record, raw []byte) (*Record) {
    return ParseJSONHelper(src, raw, sr.Type);
}

func(sr *SQL) Create(r *Request, src *Record) {
    lp, psql := sr.GetPromise(r);
    defer lp.Release();
    r.API.Logger.Debugf("CREATE GOT CONTEXT: %#v\n", psql);
    if(src.Id != "") {
        panic(errors.New("SQL does not support specifying an ID for Create() requests.")); // TODO: it should
    }
    tx, err := psql.GetSQLTransaction(sr.DB)
    if err != nil {
        panic(err);
    }
    SetId(src.Attributes, src.Id);
    err = meddler.Insert(tx, sr.Table, src.Attributes)
    if err != nil {
        panic(err);
    }
}

func(sr *SQL) Update(r *Request, rec *Record) {
    lp, psql := sr.GetPromise(r);
    defer lp.Release();
    if rec.Attributes != nil {
        SetId(rec.Attributes, rec.Id);
    } else {
        // TODO: should this panic? is it possible to UPDATE with a nil ID?
    }
    tx, err := psql.GetSQLTransaction(sr.DB);
    if err != nil {
        panic(err);
    }
    r.API.Logger.Debugf("Fields: %#v\n", rec.Attributes);
    err = meddler.Update(tx, sr.Table, rec.Attributes);
    if err != nil {
        panic(err);
    }
}

func (sr *SQL) ConvertInterfaceSliceToRecordSlice(src interface{}) []*Record {
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

func (sr *SQL) GetTableFieldFromStructField(structstr string) (string, error) {
    field, found := sr.Type.FieldByName(structstr);
    if(!found) {
        return "", errors.New("Field \""+structstr+"\" does not exist on "+sr.Type.Name());
    }
    realname := field.Name;

    meddler_tags := strings.Split(field.Tag.Get("meddler"),",");

    if(len(meddler_tags) > 0 && meddler_tags[0] != "") {
        realname = meddler_tags[0];
    }

    return realname, nil;
}
