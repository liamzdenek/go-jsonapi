package oauth2

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/RangelReale/osin/example"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type OAuth2 struct {
	Router                   *httprouter.Router
	FuncAuthorize, FuncToken httprouter.Handle
	BaseURI                  string
}

func NewOAuth2(base string) *OAuth2 {
	cfg := osin.NewServerConfig()
	cfg.AllowGetAccessRequest = true

	server := osin.NewServer(cfg, example.NewTestStorage())

	funcauthorize := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
			if !example.HandleLoginPage(ar, w, r) {
				return
			}
			ar.Authorized = true
			server.FinishAuthorizeRequest(resp, r, ar)
		}
		if resp.IsError && resp.InternalError != nil {
			fmt.Printf("ERROR: %s\n", resp.InternalError)
		}
		osin.OutputJSON(resp, w, r)
	}

	functoken := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAccessRequest(resp, r); ar != nil {
			ar.Authorized = true
			server.FinishAccessRequest(resp, r, ar)
		}
		if resp.IsError && resp.InternalError != nil {
			fmt.Printf("ERROR: %s\n", resp.InternalError)
		}
		osin.OutputJSON(resp, w, r)

	}

	o := &OAuth2{
		FuncAuthorize: funcauthorize,
		FuncToken:     functoken,
		Router:        httprouter.New(),
		BaseURI:       base,
	}
	o.InitRouter()
	return o
}

func (o *OAuth2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("OAuth2 ServeHTTP %s\n", o.BaseURI)
	o.Router.ServeHTTP(w, r)
}

func (o *OAuth2) InitRouter() {
	o.Router.GET(o.BaseURI+"authorize", o.FuncAuthorize)
	o.Router.GET(o.BaseURI+"token", o.FuncToken)
}
