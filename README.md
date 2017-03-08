# crestron-control-microservice
[![CircleCI](https://img.shields.io/circleci/project/byuoitav/crestron-control-microservice.svg)](https://circleci.com/gh/byuoitav/crestron-control-microservice) [![Apache 2 License](https://img.shields.io/hexpm/l/plug.svg)](https://raw.githubusercontent.com/byuoitav/crestron-control-microservice/master/LICENSE)

## Configuration File
The configuration file, `signal-configuration.json`, exists to make mapping of symbol names to Go code and behavior easier. In general, the configuration file is structured as follows:
```
{
    "SetVolume": {
        "signalName": "setVolumeLevel",
        "signalValue": "level",
		"parameterized": true,
        "highLow": false
    },
    "PowerOn": {
        "signalName": "confirmStartupPress",
        "signalValue": 1,
		"parameterized": true,
        "highLow": true 
    }
}
```
