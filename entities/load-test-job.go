package entities

// LoadTestJob Entity
type LoadTestJob struct {
	Name           string                 `json:"name" yaml:"name"`
	Type           string                 `json:"type" yaml:"type"`
	Target         LoadTestJobTarget      `json:"target" yaml:"target"`
	ConstantLoad   *LoadTestConstantJob   `json:"constantLoad" yaml:"constantLoad"`
	IncreasingLoad *LoadTestIncreasingJob `json:"increasingLoad" yaml:"increasingLoad"`
	FuzzyLoad      *LoadTestFuzzyJob      `json:"fuzzyLoad" yaml:"fuzzyLoad"`
}
