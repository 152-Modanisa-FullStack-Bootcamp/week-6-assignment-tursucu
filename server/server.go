package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"wallet/handler"
	"wallet/model"
	"wallet/repository"
	routes2 "wallet/routes"
	"wallet/service"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) StartServer(port int) error {
	a := http.HandlerFunc(Serve)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), a)
	return err
}

var (
	wallets       = make(map[string]*model.WalletData)
	walletRepo    = repository.NewWalletRepo(wallets)
	walletService = service.NewWalletService(walletRepo)
	walletHandler = handler.NewWalletHandler(walletService)
	routes        = []routes2.Route{
		routes2.NewRoute("GET", "/", walletHandler.Wallets),
		routes2.NewRoute("GET", "/([^/]+)", walletHandler.WalletByUsername),
		routes2.NewRoute("PUT", "/([^/]+)", walletHandler.CreateWallet),
		routes2.NewRoute("POST", "/([^/]+)", walletHandler.UpdateWalletByUsername),
	}
)

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.Regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.Method {
				allow = append(allow, route.Method)
				continue
			}
			ctx := context.WithValue(r.Context(), routes2.CtxKey{}, matches[1:])
			route.Handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}
