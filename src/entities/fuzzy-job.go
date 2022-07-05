package entities

// FuzzyJob Entity
type FuzzyJob struct {
	Duration int             `json:"duration" yaml:"duration"`
	Options  FuzzyJobOptions `json:"specs" yaml:"specs"`
}
