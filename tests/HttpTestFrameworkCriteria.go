package jsonapitest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

type CrURL struct {
	Url string
}

func SetURL(url string) *CrURL {
	return &CrURL{Url: url}
}

func (crurl *CrURL) Setup(server *httptest.Server, req *http.Request) {
	var err error
	req.URL, err = url.Parse(server.URL + crurl.Url)
	Check(err)
}

func (url *CrURL) Check(req *http.Request, res *http.Response) {

}

func (url *CrURL) Describe() string { return fmt.Sprintf("URL = %s", url.Url) }

type CrStatusCode struct {
	Code int
}

func SetStatusCode(code int) *CrStatusCode {
	return &CrStatusCode{Code: code}
}

func (crstatus *CrStatusCode) Setup(server *httptest.Server, request *http.Request) {
}

func (crstatus *CrStatusCode) Check(request *http.Request, response *http.Response) {
	if response.StatusCode != crstatus.Code {
		panic(fmt.Sprintf("Expected code %d, got %d", crstatus.Code, response.StatusCode))
	}
}

func (crstatus *CrStatusCode) Describe() string {
	return fmt.Sprintf("ExpectedStatus = %d", crstatus.Code)
}
