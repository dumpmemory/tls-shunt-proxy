package http2

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyPassHandler struct {
	proxy *httputil.ReverseProxy
}

func NewProxyPassHandler(target string) *ProxyPassHandler {
	handler := ProxyPassHandler{}

	targetUrl, err := url.Parse(target)
	if err != nil {
		log.Fatalf("parse target url %s failed", target)
	}

	if targetUrl.Scheme != "h2c" {
		log.Fatalln("http2 proxy pass supports h2c only")
	}

	targetUrl.Scheme = "https"
	handler.proxy = httputil.NewSingleHostReverseProxy(targetUrl)
	handler.proxy.Transport = &http2.Transport{
		DialTLS: func(network, addr string, cfg *tls.Config) (conn net.Conn, e error) {
			return net.Dial(network, addr)
		},
	}
	return &handler
}

func (p ProxyPassHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	p.proxy.ServeHTTP(writer, request)
}
