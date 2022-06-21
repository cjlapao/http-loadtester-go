package entities

// LoadTestConstantJobOptions Entity
type LoadTestConstantJobOptions struct {
	BlockType       string `json:"type" yaml:"type"`
	BlockInterval   int    `json:"blockInterval" yaml:"blockInterval"`
	MaxTaskInterval int    `json:"maxTaskInterval" yaml:"maxTaskInterval"`
	MinTaskInterval int    `json:"minTaskInterval" yaml:"minTaskInterval"`
	CallsPerBlock   int    `json:"callsPerBlock" yaml:"callsPerBlock"`
}
