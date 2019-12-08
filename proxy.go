package main;

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sort"
)

type Proxy struct {
  port string
  log  *Log
  proxyLogger *ProxyLogger
}

func NewProxy(log *Log) *Proxy {
	return &Proxy{
		log: log,
		proxyLogger:newProxyLogger()}
}

type ProxyHandler struct {
	proxy  *httputil.ReverseProxy
	proxyLogger *ProxyLogger
}

func (p *ProxyHandler) proxyRequest(w http.ResponseWriter, r *http.Request) {
	p.proxy.Transport = &transport{http.DefaultTransport, p.proxyLogger}
	p.proxy.ServeHTTP(w, r)
}

type transport struct {
	http.RoundTripper
	Proxy *ProxyLogger
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
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


type LogsHandler struct {
	Proxy *ProxyLogger
}

func (m *LogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path, _ := r.URL.Query()["path"]
	sortParam, _ := r.URL.Query()["sort"]
	jsonString := "not found"

	if len(path)>0 {
		requests := m.Proxy.getRegex(path[0])
		if len(requests)>0 && len(sortParam) > 0 {
			sort.Slice(requests, func(i, j int) bool {
				if sortParam[0] == "desc" {
					return requests[i].(*Request).Timestamp > requests[j].(*Request).Timestamp
				} else {
					return requests[i].(*Request).Timestamp > requests[j].(*Request).Timestamp
				}
			})
		}
		jsonString, err := json.Marshal(m.Proxy.getRegex(path[0]))
		if err != nil {}
		w.Write([]byte(jsonString))
		return
	}

	w.Write([]byte(jsonString))
}

func (p *Proxy) serveLogs() {
	log.Fatal(http.ListenAndServe(":8001", &LogsHandler{p.proxyLogger}))
}

func (p *Proxy) startProxy(port string, url *url.URL) {
	go func() {
		proxy := &ProxyHandler{proxy: httputil.NewSingleHostReverseProxy(url),
			proxyLogger:p.proxyLogger}
		http.HandleFunc("/", proxy.proxyRequest)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()
}
