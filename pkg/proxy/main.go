package proxy

import (
	"context"
	"net/http"
	"regexp"
)

type Proxy struct {
	server *http.Server
}

func NewProxy(listenAddress string, handlerMapping map[*regexp.Regexp]Handler) *Proxy {
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

func (p *Proxy) Shutdown(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

func (p *Proxy) Addr() string {
	return p.server.Addr
}
