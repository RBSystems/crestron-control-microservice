## Configuration File
The configuration file, `signal-configuration.json`, exists to make mapping of symbol names to Go code and behavior easier. In general, the configuration file is structured as follows:
```
{
    "SetVolume": {
        "signalName": "setVolumeLevel",
        "signalValue": "level",
        "parametrizable": true,
        "highLow": false
    },
    "PowerOn": {
        "signalName": "confirmStartupPress",
        "signalValue": "",
        "parametrizable": false,
        "highLow": true 
    }
}
```
