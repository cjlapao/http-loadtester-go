package entities

// LoadTestFuzzyJobOptions Entity
type LoadTestFuzzyJobOptions struct {
	BlockType        string `json:"type" yaml:"type"`
	MaxBlockInterval int    `json:"maxBlockInterval" yaml:"maxBlockInterval"`
	MinBlockInterval int    `json:"minBlockInterval" yaml:"minBlockInterval"`
	MaxTasksPerBlock int    `json:"maxTasksPerBlock" yaml:"maxTasksPerBlock"`
	MinTasksPerBlock int    `json:"minTasksPerBlock" yaml:"minTasksPerBlock"`
	MaxTaskInterval  int    `json:"maxTaskInterval" yaml:"maxTaskInterval"`
	MinTaskInterval  int    `json:"minTaskInterval" yaml:"minTaskInterval"`
}
