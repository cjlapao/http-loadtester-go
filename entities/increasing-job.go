package entities

// IncreasingJob Entity
type IncreasingJob struct {
	Duration int                  `json:"duration" yaml:"duration"`
	Options  IncreasingJobOptions `json:"specs" yaml:"specs"`
}
