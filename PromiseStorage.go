package jsonapi

import (
	"reflect"
)

type PromiseStorage struct {
	Promises map[reflect.Type]chan PromiseStorageLease
	ChanGet  chan PromiseStorageLease
}

type PromiseStorageLease struct {
	Type         reflect.Type
	Initialize   func() Promise
	ChanResponse chan LeasedPromise
}

type LeasedPromise struct {
	Promise
	ChanRelease chan bool
}

func (sp *LeasedPromise) Release() {
	close(sp.ChanRelease)
}

func NewPromiseStorage() *PromiseStorage {
	ps := &PromiseStorage{
		Promises: make(map[reflect.Type]chan PromiseStorageLease),
		ChanGet:  make(chan PromiseStorageLease),
	}
	ps.Worker()
	return ps
}

func (ps *PromiseStorage) Defer() {
	close(ps.ChanGet)
}

func (ps *PromiseStorage) Get(typ Promise, init func() Promise) LeasedPromise {
	lease := PromiseStorageLease{
		Type:         reflect.TypeOf(typ),
		Initialize:   init,
		ChanResponse: make(chan LeasedPromise),
	}
	ps.ChanGet <- lease
	return <-lease.ChanResponse
}

func (ps *PromiseStorage) Worker() {
	go func() {
	OUTER:
		for {
			select {
			case p, chanok := <-ps.ChanGet:
				if !chanok {
					break OUTER
				}
				req, ok := ps.Promises[p.Type]
				if !ok {
					req = ps.PromiseWorker(p.Initialize())
					ps.Promises[p.Type] = req
					defer close(req)
				}
				req <- p
			}
		}
	}()
}

func (ps *PromiseStorage) PromiseWorker(p Promise) chan PromiseStorageLease {
	leasechan := make(chan PromiseStorageLease)
	go func() {
		for leasereq := range leasechan {
			leased := LeasedPromise{
				Promise:     p,
				ChanRelease: make(chan bool),
			}
			leasereq.ChanResponse <- leased
			<-leased.ChanRelease
		}
	}()
	return leasechan
}
