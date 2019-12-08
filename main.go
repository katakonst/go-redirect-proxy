package main

import (
	"fmt"
	"log"
	"net/url"
)

func main() {
	appConfigs, err := InitConfig()
	if err != nil {
		log.Fatalf("Failed to load configs: %s", err)
	}

	logger := NewLogger(appConfigs.LogLevel)
	proxy := NewProxy(logger)
	if appConfigs.source !="" {
		url, err := url.Parse("http://localhost:"+appConfigs.target)
		if err != nil {
			panic(err)
		}
		proxy.startProxy(appConfigs.source, url)
		proxy.serveLogs()
	}

	hosts := appConfigs.ProxyConfigs["rules"]
	rules:= hosts.([]interface {})
	for _,elem := range rules {
		 url, err := url.Parse("http://localhost:"+
		 	fmt.Sprintf("%v", elem.(map[string]interface{})["target"]))
		if err != nil {
			panic(err)
		}
		fmt.Println(elem.(map[string]interface{})["target"])
		proxy.startProxy(fmt.Sprintf("%v", elem.(map[string]interface{})["source"]), url)
	}
	proxy.serveLogs()
}