package entities

// LoadTestIncreasingJob Entity
type LoadTestIncreasingJob struct {
	Duration int                          `json:"duration" yaml:"duration"`
	Options  LoadTestIncreasingJobOptions `json:"specs" yaml:"specs"`
}
