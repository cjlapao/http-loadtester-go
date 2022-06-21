package entities

// Job Entity
type Job struct {
	Name           string         `json:"name" yaml:"name"`
	Type           string         `json:"type" yaml:"type"`
	Target         JobTarget      `json:"target" yaml:"target"`
	ConstantLoad   *ConstantJob   `json:"constantLoad" yaml:"constantLoad"`
	IncreasingLoad *IncreasingJob `json:"increasingLoad" yaml:"increasingLoad"`
	FuzzyLoad      *FuzzyJob      `json:"fuzzyLoad" yaml:"fuzzyLoad"`
}
