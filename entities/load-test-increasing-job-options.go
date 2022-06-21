package entities

// LoadTestIncreasingJobOptions Entity
type LoadTestIncreasingJobOptions struct {
	BlockType       string `json:"type" yaml:"type"`
	BlockInterval   int    `json:"blockInterval" yaml:"blockInterval"`
	TotalCalls      int    `json:"totalCalls" yaml:"totalCalls"`
	MaxTaskInterval int    `json:"maxTaskInterval" yaml:"maxTaskInterval"`
	MinTaskInterval int    `json:"minTaskInterval" yaml:"minTaskInterval"`
}
