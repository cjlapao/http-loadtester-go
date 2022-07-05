package entities

// ConstantJob entity
type ConstantJob struct {
	Duration int                `json:"duration" yaml:"duration"`
	Options  ConstantJobOptions `json:"specs" yaml:"specs"`
}
