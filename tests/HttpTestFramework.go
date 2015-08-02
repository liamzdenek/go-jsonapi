package jsonapitest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

type TestSuite struct {
	T     *testing.T
	Tests []Test
}

type Test struct {
	Criteria []Criterion
}

type Criterion interface {
	Describe() string
	Setup(*httptest.Server, *http.Request)
	Check(*http.Request, *http.Response)
}

func NewTestSuite(T *testing.T) *TestSuite {
	return &TestSuite{
		T:     T,
		Tests: []Test{},
	}
}

func (ts *TestSuite) PushTests(tests ...Test) {
	ts.Tests = append(ts.Tests, tests...)
}

func (ts *TestSuite) PushNewTest(cr ...Criterion) {
	ts.PushTests(NewTest(cr...))
}

func (ts *TestSuite) Run(server *httptest.Server) {
	for _, test := range ts.Tests {
		req := &http.Request{}
		res := &http.Response{}
		work := func() {
			var err error
			for _, cr := range test.Criteria {
				cr.Setup(server, req)
			}

			cl := http.Client{}
			res, err = cl.Do(req)
			Check(err)

			for _, cr := range test.Criteria {
				cr.Check(req, res)
			}
		}
		err := Catch(work)
		if err != nil {
			rawReq, terr := httputil.DumpRequestOut(req, true)
			Check(terr)
			rawRes, terr := httputil.DumpResponse(res, true)
			Check(terr)
			ts.T.Errorf("Error running test: %#v\n\nTEST DESCRIPTION:\n\n%s\nREQUEST:\n\n%s\nRESPONSE:\n\n%s\n", err, test.Describe(), string(rawReq), string(rawRes))
		}
	}
}

func NewTest(c ...Criterion) Test {
	return Test{
		Criteria: c,
	}
}

func (test *Test) Describe() string {
	desc := ""
	for _, cr := range test.Criteria {
		desc = fmt.Sprintf("%s%s\n", desc, cr.Describe())
	}
	return desc
}

func Catch(f func()) (res interface{}) {
	defer func() {
		if r := recover(); r != nil {
			res = r
		}
	}()
	f()
	return
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
