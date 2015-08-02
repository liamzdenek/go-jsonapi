package jsonapitest

import (
	"fmt"
	"net/url"
)

type CrURL struct {
	Url string
}

func SetURL(url string) *CrURL {
	return &CrURL{Url: url}
}

func (crurl *CrURL) Setup(session *TestSession) {
	var err error
	session.Request.URL, err = url.Parse(session.Server.URL + crurl.Url)
	Check(err)
}

func (url *CrURL) Check(session *TestSession) {}

func (url *CrURL) Describe() string { return fmt.Sprintf("URL = %s", url.Url) }

type CrStatusCode struct {
	Code int
}

func SetStatusCode(code int) *CrStatusCode {
	return &CrStatusCode{Code: code}
}

func (crstatus *CrStatusCode) Setup(session *TestSession) {}

func (crstatus *CrStatusCode) Check(session *TestSession) {
	if session.Response.StatusCode != crstatus.Code {
		panic(fmt.Sprintf("Expected code %d, got %d", crstatus.Code, session.Response.StatusCode))
	}
}

func (crstatus *CrStatusCode) Describe() string {
	return fmt.Sprintf("ExpectedStatus = %d", crstatus.Code)
}
