package entities

// ConstantJobOptions Entity
type ConstantJobOptions struct {
	NumberOfBlocks  int    `json:"numberOfBlocks" yaml:"numberOfBlocks"`
	BlockType       string `json:"type" yaml:"type"`
	BlockInterval   int    `json:"blockInterval" yaml:"blockInterval"`
	MaxTaskInterval int    `json:"maxTaskInterval" yaml:"maxTaskInterval"`
	MinTaskInterval int    `json:"minTaskInterval" yaml:"minTaskInterval"`
	CallsPerBlock   int    `json:"callsPerBlock" yaml:"callsPerBlock"`
}
