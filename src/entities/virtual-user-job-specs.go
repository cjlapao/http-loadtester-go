package entities

// VirtualUserJobSpecs Entity
type VirtualUserJobSpecs struct {
	NumberOfVirtualUsers   int `json:"numberOfVirtualUsers" yaml:"numberOfVirtualUsers"`
	Duration               int `json:"duration" yaml:"duration"`
	MaxVirtualUserInterval int `json:"maxVirtualUserInterval" yaml:"maxVirtualUserInterval"`
	MinVirtualUserInterval int `json:"minVirtualUserInterval" yaml:"minVirtualUserInterval"`
	MaxTaskInterval        int `json:"maxTaskInterval" yaml:"maxTaskInterval"`
	MinTaskInterval        int `json:"minTaskInterval" yaml:"minTaskInterval"`
}
