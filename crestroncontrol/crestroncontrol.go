package crestroncontrol

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ParseConfig() (AllSignalConfig, error) {
	config := AllSignalConfig{}

	bytes, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/github.com/byuoitav/crestron-control-microservice/signal-configuration.json")
	if err != nil {
		return AllSignalConfig{}, err
	}

	err = json.Unmarshal(bytes, &config)
	return config, err
}
