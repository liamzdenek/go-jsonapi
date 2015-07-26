package jsonapi

import "fmt"

type FutureOutputDelete struct{}

func (fo *FutureOutputDelete) ShouldCombine(f Future) bool { return false }
func (fo *FutureOutputDelete) Combine(f Future) error      { panic("This should never be called") }

func (fo *FutureOutputDelete) Work(pf *ExecutableFuture) {
	rawreq := pf.GetRequest()
	fmt.Printf("FUTUREOUTPUT GOT RES: %#v\n", rawreq.Kind)
	switch req := rawreq.Kind.(type) {
	case *FutureRequestKindIdentity:
		fmt.Printf("FUTUREOUTPUT GOT RES: %#v\n", req.Response)
	default:
		panic(TODO())
	}
}
