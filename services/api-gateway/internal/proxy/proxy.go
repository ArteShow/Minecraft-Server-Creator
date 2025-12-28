package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(target string, rewritePath string) *httputil.ReverseProxy {
	u, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(u)

	proxy.Director = func(r *http.Request) {
		r.URL.Scheme = u.Scheme
		r.URL.Host = u.Host
		r.Host = u.Host
		r.URL.Path = rewritePath
	}

	return proxy
}