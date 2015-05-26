package jsonapie;

import (
    . ".."
    "crypto/sha256"
    "hash"
    "errors"
    "fmt"
)

type ResourceCache struct{
    Resource, Cache Resource
    CacheName string
    HashFunc func() hash.Hash
}

type CacheRecord struct {
    Id string `jsonapi:"id"`
    Iders []Ider
}

func init() {
    // safety check to make sure ResourceCache is a Resource
    var t Resource;
    t = &ResourceCache{};
    _ = t;
}

func NewResourceCache(name string, resource, cache Resource) *ResourceCache {
    return &ResourceCache{
        Resource: resource,
        Cache: cache,
        CacheName: name,
        HashFunc: sha256.New,
    }
}

func(rc *ResourceCache) GenCacheKey(args ...string) (res string) {
    h := rc.HashFunc();
    for _,arg := range args {
        h.Write([]byte(arg));
    }
    r := []byte{}
    r = h.Sum(r);
    if len(r) == 0 {
        panic("Cache hash returned by HashFunc is not populating res");
    }
    res = fmt.Sprintf("%X", r);
    return
}

func(rc *ResourceCache) CacheFind(s Session, hkey string) (res *CacheRecord, err error) {
    a := s.GetData().API;
    r := Catch(func() {
        ider, err := rc.Cache.FindOne(s, hkey);
        Check(err);
        a.Logger.Printf("CACHE OUTPUT: %v\n", ider);
        res = ider.(*CacheRecord);
    });
    if r != nil {
        err = errors.New(fmt.Sprintf("Caught error from Cache: %#v\n", r));
    }
    return;
}

func(rc *ResourceCache) CacheCreate(s Session, hkey string, r *CacheRecord) {
    a := s.GetData().API;
    a.Logger.Printf("PUTTING INTO CACHE: %s - %#v\n", hkey, r);
    rc.Cache.Create(s,r,&hkey)
}

func(rc *ResourceCache) FindDefault(s Session, rp RequestParams) ([]Ider, error) {
    panic("TODO");
}

func(rc *ResourceCache) FindOne(s Session, id string) (Ider, error) {
    a := s.GetData().API;
    hkey := string(rc.GenCacheKey("FindOne", id));
    a.Logger.Printf("CHECKING CACHE\n");
    cacherecord, err := rc.CacheFind(s,hkey);
    if err == nil && cacherecord != nil {
        a.Logger.Printf("CACHE SUCCESS %#v\n", cacherecord);
        return cacherecord.Iders[0], nil;
    }
    a.Logger.Printf("CACHE FAILURE -- %s\n", err);
    ider, err := rc.Resource.FindOne(s,id);
    if err == nil && ider != nil {
        rc.CacheCreate(s,hkey, &CacheRecord{
            Iders: []Ider{ider},
        });
    }
    return ider, err;
}

func(rc *ResourceCache) FindMany(s Session, rp RequestParams, ids []string) ([]Ider, error) {
    return rc.Resource.FindMany(s,rp, ids);
}

// TODO: this function should be cached
func(rc *ResourceCache) FindManyByField(s Session, rp RequestParams, field string, value string) ([]Ider, error) {
    return rc.Resource.FindManyByField(s,rp, field, value);
}

func(rc *ResourceCache) Delete(s Session, id string) error {
    return rc.Resource.Delete(s,id);
}

func(rc *ResourceCache) ParseJSON(s Session, idersrc Ider, raw []byte) (Ider, *string, *string, *OutputLinkageSet, error) {
    return rc.Resource.ParseJSON(s, idersrc, raw);
}

func(rc *ResourceCache) Create(s Session, ider Ider, id *string) (RecordCreatedStatus, error) {
    return rc.Resource.Create(s, ider, id);
}

func(rc *ResourceCache) Update(s Session, id string, ider Ider) error {
    panic("NOT IMPLEMENTED");
}
