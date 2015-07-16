package resource;

import (
    "database/sql";
    "github.com/russross/meddler";
    "reflect"
    "strings"
    "fmt"
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

func(f *FutureSQL) RunQuery(pf *ExecutableFuture, req *FutureRequest, parameters []SQLExpression) ([]*Record, bool, *OError){
    lp, psql := f.Resource.GetPromise(pf.Request);
    defer lp.Release();

    tx, err := psql.GetSQLTransaction(f.Resource.DB)
    if err != nil {
        oe := ErrorToOError(err)
        return []*Record{}, false, &oe;
    }
    vs := reflect.New(reflect.SliceOf(reflect.PtrTo(f.Resource.Type))).Interface()
    query, queryargs, is_single := f.PrepareQuery(parameters...);
    pf.Request.API.Logger.Debugf("RUN QUERY: %#v %#v\n", query, queryargs);
    err = meddler.QueryAll(
        tx,
        vs,
        query,
        queryargs...,
    );
    var oerr *OError;
    if err != nil {
        oerr = &OError{
            Title: err.Error(),
            Detail: fmt.Sprintf("%s -- %#v\n", query, queryargs),
        };
    }
    return ConvertInterfaceSliceToRecordSlice(vs), is_single, oerr;
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
            Failure: []OError{*err},
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
            Failure: []OError{*err},
        });
        return;
    }
    records_by_field := map[Field][]*Record{};
    for _, record := range records {
        key := Field{Field: f.Resource.GetIdFieldName(nil), Value: record.Id};
        records_by_field[key] = append(records_by_field[key], record);
    }
    req.SendResponse(&FutureResponse{
        IsSuccess: true,
        Success: map[Future]FutureResponseKind{
            f: &FutureResponseKindByFields{
                IsSingle: forced_single || is_single,
                Records: records_by_field,
                //Field{Field: f.Resource.GetIdFieldName(nil), Value:k.Ids[0]}: records,

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

func(sr *SQL) GetFuture() Future {
    return &FutureSQL{
        Resource: sr,
    }
}
/*
func(sr *SQL) ParseJSON(r *Request, src *Record, raw []byte) (*Record) {
    return ParseJSONHelper(src, raw, sr.Type);
}
*/
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
/*

func(sr *SQL) Delete(r *Request, id string) {
    id_sql_name := sr.GetIdFieldName(nil);
    _, err := sr.DB.Exec("DELETE FROM "+sr.Table+" WHERE "+id_sql_name+"=?", id);
    if err != nil {
        panic(err);
    }
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
*/
