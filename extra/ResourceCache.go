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
    res = fmt.Sprintf("val: %X", r);
    return
}

func(rc *ResourceCache) CacheFind(a *API, s Session, hkey string) (res *CacheRecord, err error) {
    r := Catch(func() {
        ider, err := rc.Cache.FindOne(a, s, hkey);
        Check(err);
        fmt.Printf("CACHE OUTPUT: %v\n", ider);
        res = ider.(*CacheRecord);
    });
    if r != nil {
        err = errors.New(fmt.Sprintf("Caught error from Cache: %#v\n", r));
    }
    return;
}

func(rc *ResourceCache) CacheCreate(a *API, s Session, hkey string, r *CacheRecord) {
    fmt.Printf("PUTTING INTO CACHE: %s - %#v\n", hkey, r);
    rc.Cache.Create(a,s,r,&hkey)
}

func(rc *ResourceCache) FindOne(a *API, s Session, id string) (Ider, error) {
    hkey := string(rc.GenCacheKey("FindOne", id));
    fmt.Printf("CHECKING CACHE\n");
    cacherecord, err := rc.CacheFind(a,s,hkey);
    if err == nil && cacherecord != nil {
        fmt.Printf("CACHE SUCCESS %#v\n", cacherecord);
        return cacherecord.Iders[0], nil;
    }
    fmt.Printf("CACHE FAILURE -- %s\n", err);
    ider, err := rc.Resource.FindOne(a,s,id);
    if err == nil && ider != nil {
        rc.CacheCreate(a,s,hkey, &CacheRecord{
            Iders: []Ider{ider},
        });
    }
    return ider, err;
}

func(rc *ResourceCache) FindMany(a *API, s Session, p *Paginator, ids []string) ([]Ider, error) {
    return rc.Resource.FindMany(a, s,p, ids);
}

func(rc *ResourceCache) FindManyByField(a *API, s Session, field string, value string) ([]Ider, error) {
    return rc.Resource.FindManyByField(a,s, field, value);
}

func(rc *ResourceCache) Delete(a *API, s Session, id string) error {
    return rc.Resource.Delete(a,s,id);
}

func(rc *ResourceCache) ParseJSON(a *API, s Session, raw []byte) (Ider, *string, *string, *OutputLinkageSet, error) {
    return rc.Resource.ParseJSON(a, s, raw);
}

func(rc *ResourceCache) Create(a *API, s Session, ider Ider, id *string) (RecordCreatedStatus, error) {
    return rc.Resource.Create(a,s, ider, id);
}
