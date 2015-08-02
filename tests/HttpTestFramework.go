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

type TestSession struct {
	Test     *Test
	Server   *httptest.Server
	Request  *http.Request
	Response *http.Response
}

type Criterion interface {
	Describe() string
	Setup(*TestSession)
	Check(*TestSession)
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
		session := &TestSession{
			Server:  server,
			Test:    &test,
			Request: &http.Request{},
		}
		work := func() {
			for _, cr := range test.Criteria {
				cr.Setup(session)
			}

			cl := http.Client{}
			res, err := cl.Do(session.Request)
			Check(err)

			session.Response = res

			for _, cr := range test.Criteria {
				cr.Check(session)
			}
		}
		err := Catch(work)
		if err != nil {
			rawReq, terr := httputil.DumpRequestOut(session.Request, true)
			Check(terr)
			rawRes, terr := httputil.DumpResponse(session.Response, true)
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
