package transport

import (
	"bioly/gateway/internal/config"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Handler struct {
	proxies map[string]*httputil.ReverseProxy
}

func NewHandler(proxies []config.Proxy) (*Handler, error) {
	if len(proxies) == 0 {
		return nil, fmt.Errorf("got empty proxies list, nothing to handle")
	}

	handlerProxies := make(map[string]*httputil.ReverseProxy, len(proxies))
	for _, proxy := range proxies {
		u, err := url.Parse(proxy.Url)
		if err != nil {
			return nil, err
		}

		handlerProxies[proxy.Endpoint] = httputil.NewSingleHostReverseProxy(u)
	}

	return &Handler{proxies: handlerProxies}, nil
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	for svc, proxy := range h.proxies {
		prefix := "/" + svc + "/"
		mux.Handle(prefix, http.StripPrefix(prefix, proxy))
	}
}
