package main

import (
	"bytes"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type Proxy struct {
  port string
  log  *Log
  proxyLogger *ProxyLogger
  logPort *string
}

func NewProxy(log *Log, logPort *string) *Proxy {
	return &Proxy{
		log: log,
		proxyLogger:newProxyLogger(log),
		logPort:logPort}
}

type ProxyHandler struct {
	proxy  *httputil.ReverseProxy
	proxyLogger *ProxyLogger
}

func (p *ProxyHandler) proxyRequest(w http.ResponseWriter, r *http.Request) {
	p.proxy.Transport = &Transport{http.DefaultTransport, p.proxyLogger}
	p.proxy.ServeHTTP(w, r)
}

func (p *Proxy) serveLogs() {
	log.Fatal(http.ListenAndServe(":" + *p.logPort, &LogsHandler{p.proxyLogger,p.log}))
}

func (p *Proxy) startProxy(port string, url *url.URL) {
	go func() {
		proxy := &ProxyHandler{proxy: httputil.NewSingleHostReverseProxy(url),
			proxyLogger:p.proxyLogger}
		http.HandleFunc("/", proxy.proxyRequest)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()
}

type Transport struct {
	http.RoundTripper
	Proxy *ProxyLogger
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	t.Proxy.addToMap(html.EscapeString(resp.Request.URL.Path),
		string(b), resp.Request.Method, req.Header)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))

	return resp, nil
}
