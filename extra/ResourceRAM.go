package jsonapie;

import (
    "reflect"
    "strings"
    "fmt"
    "errors"
    "encoding/json"
    . ".."
);

type ResourceRAM struct{
    Type reflect.Type
    Storage map[string]Ider
}

func init() {
    // safety check to make sure ResourceRAM is a Resource
    var t Resource;
    t = &ResourceRAM{};
    _ = t;
}

func NewResourceRAM(t Ider) *ResourceRAM {
    return &ResourceRAM{
        Type: reflect.Indirect(reflect.ValueOf(t)).Type(),
        Storage: map[string]Ider{},
    }
}

func(rr *ResourceRAM) FindOne(id string) (Ider, error) {
    val, exists := rr.Storage[id];
    if(!exists) {
        return nil, nil;
    }
    return val.(Ider), nil;
}

func(rr *ResourceRAM) FindMany(ids []string) ([]Ider, error) {
    res := []Ider{};
    for _, id := range ids {
        val, err := rr.FindOne(id);
        if err == nil && val != nil {
            res = append(res, val);
        }
    }
    return res, nil
}

func(rr *ResourceRAM) FindManyByField(field string, value string) ([]Ider, error) {
    return nil, errors.New("ResourceRAM does not support FindManyByField -- you are probably using a linkage with a ResourceRAM as the target");
}

func(rr *ResourceRAM) Delete(id string) error {
    if _,exists := rr.Storage[id]; exists {
        delete(rr.Storage, id);
    }
    return nil;
}

func(rr *ResourceRAM) Create(resource_str string, raw []byte, verify ResourceCreateVerifyFunc) (Ider, RecordCreatedStatus, error) {
    v := reflect.New(rr.Type).Interface();
    rp := NewRecordParserSimple(v);
    err := json.Unmarshal(raw, rp);
    if(err != nil) {
        return nil, StatusFailed, err;
    }
    if(rp.Data.Id == nil) {
        return nil, StatusFailed, errors.New("ResourceRAM requires specifying an ID for Create() requests."); // TODO: it should
    }
    if(rp.Data.Type != resource_str) {
        return nil, StatusFailed, errors.New(fmt.Sprintf("This is resource \"%s\" but the new object includes type:\"%s\"", resource_str, rp.Data.Type));
    }
    ider := v.(Ider);
    if err = verify(ider, rp.Linkages()); err != nil {
        return nil, StatusFailed, errors.New(fmt.Sprintf("The linkage verification function returned an error: %s\n", err));
    }
    rr.Storage[ider.Id()] = ider;
    return ider, StatusCreated, nil;
}

func (rr *ResourceRAM) GetTableFieldFromStructField(structstr string) (string, error) {
    field, found := rr.Type.FieldByName(structstr);
    if(!found) {
        return "", errors.New("Field "+structstr+" does not exist on "+rr.Type.Name());
    }
    realname := field.Name;

    meddler_tags := strings.Split(field.Tag.Get("meddler"),",");

    if(len(meddler_tags) > 0 && meddler_tags[0] != "") {
        realname = meddler_tags[0];
    }

    return realname, nil;
}
