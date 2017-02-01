package crestroncontrol

/*
SignalConfig represents the structure, along with SignalValue.
JSON Format of the file is as follows:

{
    "PowerStandby": {
        "signalName": "confirm_system_off",
        "isUrlParameter": true,
        "signalValueSequence": [{
            "value": "1",
            "isUrlParameter": false
        }, {
            "value": "0",
            "isUrlParameter": false
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
	SignalName          string        `json:"signalName"`
	IsURLParameter      bool          `json:"isUrlParameter"`
	SignalValueSequence []SignalValue `json:"signalValueSequence"`
}

//SignalValue is the single value to set a signal to.
type SignalValue struct {
	Value          string `json:"value"`
	IsURLParameter bool   `json:"isUrlParameter"`
}

//SignalState represents a signal state.
type SignalState struct {
	SignalName string
	Value      string
}
