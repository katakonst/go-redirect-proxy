package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"regexp"
	"time"
)

type ProxyLogger struct {
	requestsMap map[string][]interface{}
	log *Log
	wsChanel chan interface{}
}

type Request struct {
	Path string
	Timestamp int64
	Body string
	RequestBody string
	Method string
	Headers http.Header
	Status string
}

func newProxyLogger(logger *Log) *ProxyLogger {
	requestsMap := make(map[string][]interface{})
	return &ProxyLogger{requestsMap, logger, make(chan interface{},10)}
}

func (p* ProxyLogger) addToMap(path string, body string, method string,
	headers http.Header, status string, requestBody string) {

	requestsSlice := p.requestsMap[path]
	if len(requestsSlice) == 0 {
		requestsSlice = make([]interface{}, 0)
	}
	request := Request{Path:path,
		Timestamp:time.Now().Unix(),
		Body:body,
		RequestBody: requestBody,
		Method:method,
		Headers:headers}

	requestsSlice = append(requestsSlice, request)
	p.wsChanel <-request
	p.requestsMap[path]=requestsSlice
}

func (p *ProxyLogger) startWS(port string) {
	go func() {
		server2 := http.NewServeMux()
		server2.HandleFunc("/ws", p.handleWS)
		p.log.Infof("Starting ws server on port %v", port)
		log.Fatal(http.ListenAndServe(":"+port, server2))
	}()
}

var upgrader = websocket.Upgrader{}

func (p *ProxyLogger) handleWS(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		p.log.Errorf("handleWS upgrade error: %v", err)
		return
	}
	defer c.Close()
	for {
		ele:= <- p.wsChanel
		err := c.WriteJSON(ele)
		if err != nil {
			p.log.Errorf("handleWS write error: %v", err)
			break
		}
	}
}

func (p* ProxyLogger) get(path string) []interface{} {
	return p.requestsMap[path]
}

func (p* ProxyLogger) getLogByRegex(regexPath string) []interface{} {

	result :=make([]interface{}, 0)
	for k, v := range p.requestsMap {
		 matched, err := regexp.MatchString(regexPath, k)
		 if err != nil {
		   p.log.Errorf("GetLogByRegex: Error while searching log: %v", err)
		 }
		if matched {
			for _,elem := range v {
				result = append(result, elem)
			}
		 }
	}
	return result
}