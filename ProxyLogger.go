package main

import (
	"net/http"
	"regexp"
	"time"
)

type ProxyLogger struct {
	requestsMap map[string][]interface{}
}

type Request struct {
	Path string
	Timestamp int64
	Body string
	Method string
	Headers http.Header
}

func newProxyLogger() *ProxyLogger {
	requestsMap := make(map[string][]interface{})
	return &ProxyLogger{requestsMap}
}

func (p* ProxyLogger) addToMap(path string, body string, method string, headers http.Header) {
	r := p.requestsMap[path]
	if len(r) == 0 {
		r = make([]interface{}, 0)
	}
	r = append(r, Request{Path:path,
		Timestamp:time.Now().Unix(),
		Body:body, Method:method, Headers:headers})
	p.requestsMap[path]=r
}

func (p* ProxyLogger) get(path string) []interface{} {
	return p.requestsMap[path]
}

func (p* ProxyLogger) getRegex(regexPath string) []interface{} {

	result :=make([]interface{}, 0)
	for k, v := range p.requestsMap {
		 matched, err := regexp.MatchString(regexPath, k)
		 if err != nil {

		 }
		if matched {
			for _,elem := range v {
				result = append(result, elem)
			}
		 }
	}
	return result
}


