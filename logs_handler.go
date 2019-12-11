package main

import (
	"encoding/json"
	"net/http"
	"sort"
)

type LogsHandler struct {
	Proxy *ProxyLogger
	log  *Log
}

func (l *LogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path, _ := r.URL.Query()["path"]
	sortParam, _ := r.URL.Query()["sort"]
	jsonString := "not found"

	if len(path)>0 {
		requests := l.Proxy.getLogByRegex(path[0])
		if len(requests)>0 && len(sortParam) > 0 {
			sort.Slice(requests, func(i, j int) bool {
				if sortParam[0] == "desc" {
					return requests[i].(*Request).Timestamp > requests[j].(*Request).Timestamp
				} else {
					return requests[i].(*Request).Timestamp > requests[j].(*Request).Timestamp
				}
			})
		}

		jsonString, err := json.Marshal(l.Proxy.getLogByRegex(path[0]))
		if err != nil {
			l.log.Errorf("ServeHTTP: Error while unmarshaling entity %v", err)
		}
		_, err =w.Write([]byte(jsonString))
		if err!=nil {
			l.log.Errorf("ServeHTTP: Error while serving log server %v", err)
		}
		return
	}

	_, err :=w.Write([]byte(jsonString))
	if err!=nil {
		l.log.Errorf("ServeHTTP: Error while serving log server %v", err)
	}
}