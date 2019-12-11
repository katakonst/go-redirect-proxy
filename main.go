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
	proxy := NewProxy(logger, &appConfigs.logPort)
	if appConfigs.source !="" {
		targetUrl, err := url.Parse("http://localhost:"+appConfigs.target)
		if err != nil {
			panic(err)
		}
		proxy.startProxy(appConfigs.source, targetUrl)
		proxy.serveLogs()
	}

	hosts := appConfigs.ProxyConfigs["rules"]
	rules:= hosts.([]interface {})
	for _,elem := range rules {
		 ruleUrl, err := url.Parse("http://localhost:"+
		 	fmt.Sprintf("%v", elem.(map[string]interface{})["target"]))
		if err != nil {
			logger.Fatalf("Fatal error %v", err)
		}
		proxy.startProxy(fmt.Sprintf("%v", elem.(map[string]interface{})["source"]), ruleUrl)
	}
	proxy.serveLogs()
}