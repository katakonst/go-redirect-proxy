package main;

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

type Config struct {
	ProxyConfigs      map[string]interface{}
	LogLevel        string
	source          string
	target          string
	logPort         string
	wsPort          string
}

func InitConfig() (Config, error) {
	fileName := flag.String("file", "config.json", "config filename")
	logLevel := flag.String("log-level", "info", "log level")
	source := flag.String("source", "", "source port")
	target := flag.String("target", "", "target port")
	logPort := flag.String("logPort", "8081", "log port")
	wsPort := flag.String("wsPort", "3030", "ws port")
	flag.Parse()

	proxyConfigs := make(map[string]interface{})
	if *source == "" {
		var err error
		proxyConfigs, err = parseFile(*fileName)
		if err != nil {
			return Config{}, err
		}
	}


	return Config{
		ProxyConfigs:      proxyConfigs,
		LogLevel:      *logLevel,
		source:        *source,
		target:        *target,
		logPort:       *logPort,
		wsPort:        *wsPort,
	}, nil
}

func parseFile(filePath string) (map[string]interface{}, error) {
	fileContents := make(map[string]interface{})
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &fileContents); err != nil {
		return nil, err
	}

	return fileContents, nil
}