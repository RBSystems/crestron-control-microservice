package crestroncontrol

type AllSignalConfig struct {
	Mapping map[string]SignalConfig `json: "mapping"`
}

/*
SignalConfig represents the structure, along with SignalValue.
JSON Format of the file is as follows:

{
	"PowerStandby": {
      "signalName": "confirm_system_off",
      "isUrlParameter": true,
      "valueProgression": [
        {
          "value": "1",
          "isUrlParameter": false
        },
        {
          "value": "0",
          "isUrlParameter": false
        }
      ]
}
*/
type SignalConfig struct {
	SignalName       string        `json: "signalName"`
	IsUrlParameter   bool          `json: "isUrlParameter"`
	ValueProgression []SignalValue `json: "valueProgression"`
}

type SignalValue struct {
	Value          string
	IsUrlParameter bool `json: "isUrlParameter`
}
