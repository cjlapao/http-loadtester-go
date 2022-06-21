package entities

// LoadTestJobOutput Entity
type LoadTestJobOutput struct {
	MaxTaskOutput    int    `json:"maxTaskOutput" yaml:"maxTaskOutput"`
	OutputResults    bool   `json:"outputResults" yaml:"outputResults"`
	OutputToFile     bool   `json:"outputToFile" yaml:"outputToFile"`
	OutputToFilePath string `json:"outputToFilePath" yaml:"outputToFilePath"`
}
