package crestroncontrol

type AllSignalConfig struct {
	Mapping map[string]SignalConfig `json: "mapping"`
}

type SignalConfig struct {
	SignalName     string `json: "signalName"`
	SignalValue    string `json: "signalValue"`
	Parametrizable bool   `json: "parameterizable"`
	HighLow        bool   `json: "highLow"`
}
