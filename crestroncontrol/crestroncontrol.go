package crestroncontrol

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/labstack/echo"
)

//SignalConfigFile maps signal names to their operations.
var SignalConfigFile map[string]SignalConfig

//ParseConfig gets the config file and reads it into the struct
func ParseConfig() (map[string]SignalConfig, error) {
	config := make(map[string]SignalConfig)

	bytes, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/github.com/byuoitav/crestron-control-microservice/signal-configuration.json")
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	return config, err
}

//GetSignalConfigSequence needs to handle if we need to paramterize the
//signal name, as well as the value.
//returns a map of signalName -> value
func GetSignalConfigSequence(context echo.Context, command string) []SignalState {
	toReturn := []SignalState{}

	//get the SignalConfig
	config := SignalConfigFile[command]

	//Get the signal name, if parameterized, pull from the context.
	var signalName string
	if config.IsURLParameter {
		signalName = context.Param(config.SignalName)
	} else {
		signalName = config.SignalName
	}

	//Build our progression.
	for _, value := range config.SignalValueSequence {
		if value.IsURLParameter {
			toReturn = append(toReturn, SignalState{
				SignalName: signalName,
				Value:      context.Param(value.Value),
			})
		} else {
			toReturn = append(toReturn, SignalState{
				SignalName: signalName,
				Value:      value.Value,
			})
		}
	}
	return toReturn

}
