package crestroncontrol

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/labstack/echo"
)

//SignalConfigFile maps signal names to their operations.
var SignalConfigFile AllSignalConfig

//ParseConfig gets the config file and reads it into the struct
func ParseConfig() (AllSignalConfig, error) {
	config := AllSignalConfig{}

	bytes, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/github.com/byuoitav/crestron-control-microservice/signal-configuration.json")
	if err != nil {
		return AllSignalConfig{}, err
	}

	err = json.Unmarshal(bytes, &config)
	return config, err
}

//GetSignalConfigValue needs to handle if we need to paramterize the
//signal name, as well as the value.
func GetSignalConfigValue(context echo.Context, signal string) string {
	value := SignalConfigFile.Mapping[signal].SignalValue

	if SignalConfigFile.Mapping[signal].HighLow {
		return "1"
	}

	if SignalConfigFile.Mapping[signal].Parameterized {
		value = context.Param(SignalConfigFile.Mapping[signal].SignalValue)
	}

	return value
}
