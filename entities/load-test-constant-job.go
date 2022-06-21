package entities

// LoadTestConstantJob entity
type LoadTestConstantJob struct {
	Duration int                        `json:"duration" yaml:"duration"`
	Options  LoadTestConstantJobOptions `json:"specs" yaml:"specs"`
}
