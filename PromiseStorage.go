package jsonapi;

import ("reflect";)

type PromiseStorage struct {
    Promises []Promise
    ChanPush chan Promise
    ChanGet chan PromiseStorageGet
}

type PromiseStorageGet struct {
    Type reflect.Type
    ChanResponse chan Promise
}

func NewPromiseStorage() *PromiseStorage {
    ps := &PromiseStorage{
        ChanPush: make(chan Promise),
        ChanGet: make(chan PromiseStorageGet),
    };
    ps.Worker();
    return ps;
}

func(ps *PromiseStorage) Defer() {
    close(ps.ChanPush);
    close(ps.ChanGet);
}

func(ps *PromiseStorage) Worker() {
    go func() {
        for {
            select {
            case p := <-ps.ChanPush:
                ps.Promises = append(ps.Promises, p);
            case p := <-ps.ChanGet:
                //ps.Promises
            }
        }
    }();
}
