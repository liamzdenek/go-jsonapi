package jsonapie;

import (
    "reflect"
    "strings"
    "errors"
    . ".."
);

// TODO: this resource should not return the exact same pointer every time... failure semantics mean that we might accidentally push modifications that were refused due to the same data existing both here and being returned from here
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

func(rr *ResourceRAM) FindDefault(a *API, s Session, p *Paginator) ([]Ider, error) {
    panic("ResourceRAM does not support requests to FindDefault");
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

func(rr *ResourceRAM) ParseJSON(a *API, s Session, ider Ider, raw []byte) (Ider, *string, *string, *OutputLinkageSet, error) {
    return ParseJSONHelper(ider, raw, rr.Type);
}

func(rr *ResourceRAM) Create(a *API, s Session, ider Ider, id *string) (RecordCreatedStatus, error) {
    if(id == nil) {
        return StatusFailed, errors.New("ResourceRAM requires specifying an ID for Create() requests."); // TODO: it should
    }
    if _, exists := rr.Storage[*id]; exists {
        return StatusFailed, errors.New("The provided ID already exists"); // TODO: it should
    }
    SetId(ider, *id);
    a.Logger.Printf("Setting %s = %#v\n", GetId(ider), ider);
    rr.Storage[GetId(ider)] = ider;
    return StatusCreated, nil;
}

func(rr *ResourceRAM) Update(a *API, s Session, id string, ider Ider) error {
    err := SetId(ider, id);
    if err != nil {
        return err;
    }
    rr.Storage[GetId(ider)] = ider;
    return nil
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
