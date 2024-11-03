package proxy

import (
	"context"
	"net/http"
)

type Proxy struct {
	server *http.Server
}

func NewProxy(listenAddress string, handlerMapping map[string]Handler) *Proxy {
	handler := func(w http.ResponseWriter, r *http.Request) {
		defaultHandleRequest(w, r, handlerMapping)
	}

	proxy := &Proxy{
		server: &http.Server{
			Addr:    listenAddress,
			Handler: http.HandlerFunc(handler),
		},
	}

	return proxy
}

func (p *Proxy) Start() error {
	return p.server.ListenAndServe()
}

func (p *Proxy) Shutdown() error {
	return p.server.Shutdown(context.Background())
}

func (p *Proxy) Addr() string {
	return p.server.Addr
}
