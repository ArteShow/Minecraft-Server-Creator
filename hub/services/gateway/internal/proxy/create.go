package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(target string, rewritePath string) *httputil.ReverseProxy {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = u.Scheme
			r.URL.Host = u.Host
			r.Host = u.Host
			r.URL.Path = rewritePath
		},
	}
}
