package crestroncontrol

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/labstack/echo"
)

//SignalConfigFile maps signal names to their operations.
var SignalConfigFile map[string]SignalConfig

//ParseConfig gets the config file and reads it into the struct
func ParseConfig() (map[string]SignalConfig, error) {
	config := make(map[string]SignalConfig)

	fileLocation := "/go/crestron-control-microservice/signal-configuration.json" // the location of signal-configuration.json in the ARM Docker container

	if len(os.Getenv("GOPATH")) > 0 {
		fileLocation = os.Getenv("GOPATH") + "/src/github.com/byuoitav/crestron-control-microservice/signal-configuration.json" // for non-Pi deployment/development
	}

	bytes, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	return config, err
}

//GetSignalConfigSequence needs to handle if we need to paramterize the
//signal name, as well as the value.
//returns a map of signalName -> value
func GetSignalConfigSequence(context echo.Context, command string) ([]SignalState, error) {
	toReturn := []SignalState{}

	//get the SignalConfig
	config, ok := SignalConfigFile[command]
	if !ok {
		errorString := fmt.Sprintf("ERROR: No entry in config file for %v.", command)
		log.Printf(errorString)
		return toReturn, errors.New(errorString)
	}

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
	return toReturn, nil

}
