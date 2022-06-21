package entities

// LoadTestFuzzyJob Entity
type LoadTestFuzzyJob struct {
	Duration int                     `json:"duration" yaml:"duration"`
	Options  LoadTestFuzzyJobOptions `json:"specs" yaml:"specs"`
}
