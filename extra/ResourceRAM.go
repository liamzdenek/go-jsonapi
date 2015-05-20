package jsonapie;

import (
    "reflect"
    "strings"
    "errors"
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

func(rr *ResourceRAM) Push(id string, ider Ider) {
    SetId(ider, id);
    rr.Storage[id] = ider;
}

func(rr *ResourceRAM) FindOne(a *API, s Session, id string) (Ider, error) {
    val, exists := rr.Storage[id];
    if(!exists) {
        return nil, nil;
    }
    return val.(Ider), nil;
}

func(rr *ResourceRAM) FindMany(a *API, s Session, p *Paginator, ids []string) ([]Ider, error) {
    res := []Ider{};
    for _, id := range ids {
        val, err := rr.FindOne(a,s,id);
        if err == nil && val != nil {
            res = append(res, val);
        }
    }
    return res, nil
}

func(rr *ResourceRAM) FindManyByField(a *API, s Session, field string, value string) ([]Ider, error) {
    return nil, errors.New("ResourceRAM does not support FindManyByField -- you are probably using a linkage with a ResourceRAM as the target");
}

func(rr *ResourceRAM) Delete(a *API, s Session, id string) error {
    if _,exists := rr.Storage[id]; exists {
        delete(rr.Storage, id);
    }
    return nil;
}

func(rr *ResourceRAM) ParseJSON(a *API, s Session, raw []byte) (Ider, *string, *string, *OutputLinkageSet, error) {
    return ParseJSONHelper(raw, rr.Type);
}

func(rr *ResourceRAM) Create(a *API, s Session, ider Ider, id *string) (RecordCreatedStatus, error) {
    if(id == nil) {
        return StatusFailed, errors.New("ResourceRAM requires specifying an ID for Create() requests."); // TODO: it should
    }
    SetId(ider, *id);
    a.Logger.Printf("Setting %s = %#v\n", GetId(ider), ider);
    rr.Storage[GetId(ider)] = ider;
    return StatusCreated, nil;
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
