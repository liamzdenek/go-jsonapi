package jsonapie;

import (
    . ".."
    "crypto/sha256"
    "hash"
)

type ResourceCache struct{
    Resource, Cache Resource
    CacheName string
    HashFunc func() hash.Hash
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

func(rc *ResourceCache) GenCacheKey(args ...string) (res []byte) {
    h := rc.HashFunc();
    for _,arg := range args {
        h.Write([]byte(arg));
    }
    h.Sum(res);
    return
}

func(rc *ResourceCache) Capture(f func()) (r interface{}) {
    defer func() {
        r = recover();
    }();
    f();
    return nil
}

func(rc *ResourceCache) CacheFindOne(hkey string) (Ider, error) {

}

func(rc *ResourceCache) FindOne(id string) (Ider, error) {
    hkey := rc.GenCacheKey("FindOne", id);
    rc.CacheCheck(hkey);
    return rc.Resource.FindOne(id);
}

func(rc *ResourceCache) FindMany(p *Paginator, ids []string) ([]Ider, error) {
    return rc.Resource.FindMany(p, ids);
}

func(rc *ResourceCache) FindManyByField(field string, value string) ([]Ider, error) {
    return rc.Resource.FindManyByField(field, value);
}

func(rc *ResourceCache) Delete(id string) error {
    return rc.Resource.Delete(id);
}

func(rc *ResourceCache) ParseJSON(raw []byte) (Ider, *string, *string, *OutputLinkageSet, error) {
    return rc.Resource.ParseJSON(raw);
}

func(rc *ResourceCache) Create(ctx Context, resource_str string, ider Ider, id *string) (RecordCreatedStatus, error) {
    return rc.Resource.Create(ctx, resource_str, ider, id);
}
