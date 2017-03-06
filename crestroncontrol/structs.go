package crestroncontrol

/*
SignalConfig represents the structure, along with SignalValue.
JSON Format of the file is as follows:

{
    "PowerStandby": {
       "signalValueSequence": [{
			"signalName": "confirm_system_off",
			"isSignalUrlParameter": true,
			"value": "1",
			"isValueUrlParameter": false
        }, {
			"signalName": "confirm_system_off",
			"isSignalUrlParameter": true,
			"value": "1",
			"isValueUrlParameter": false
        }]
    },
    "ChangeInput": {
        "signalName": "port",
        "isUrlParameter": true,
        "signalValueSequence": [{
            "value": "1",
            "isUrlParameter": false
        }]
    },
    "SetVolume": {
        "signalName": "system_volume_level",
        "isUrlParameter": false,
        "signalValueSequence": [{
            "value": "level",
            "isUrlParameter": true
        }]
    }
}
*/
type SignalConfig struct {
	SignalValueSequence []SignalValue `json:"signalValueSequence"`
}

//SignalValue is the single value to set a signal to.
type SignalValue struct {
	SignalName           string `json:"signalName"`
	IsSignalURLParameter bool   `json:"isSignalUrlParameter"`
	Value                string `json:"value"`
	IsValueURLParameter  bool   `json:"isValueUrlParameter"`
}

//SignalState represents a signal state.
type SignalState struct {
	SignalName string
	Value      string
}
