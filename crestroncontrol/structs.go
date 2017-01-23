package crestroncontrol

type AllSignalConfig struct {
	Mapping map[string]SignalConfig `json: "mapping"`
}

type SignalConfig struct {
	SignalName    string `json: "signalName"`
	SignalValue   string `json: "signalValue"`
	Parameterized bool   `json: "parameterized"`
	HighLow       bool   `json: "highLow"`
}
