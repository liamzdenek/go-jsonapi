package jsonapitest

import (
	".."
	"../resource"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/http/httptest"
	"testing"
)

type User struct {
	Id   string `jsonapi:"id"`
	Name string
}

func WithTestServer(f func(*httptest.Server)) {
	api := jsonapi.NewAPI("/")

	db, err := sql.Open("mysql", "tasky:tasky@/tasky")

	if err != nil {
		panic(err)
	}

	resource_user := resource.NewSQL(db, "users", &User{})

	api.MountResource("user", resource_user)

	server := httptest.NewServer(api)
	defer server.Close()
	f(server)
}

func TestBasicRead(t *testing.T) {
	suite := NewTestSuite(t)

	suite.PushNewTest(
		SetURL("/user/1"),
		SetStatusCode(200),
	)

	WithTestServer(func(server *httptest.Server) {
		suite.Run(server)
	})
}
