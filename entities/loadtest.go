package entities

type LoadTest struct {
	DisplayName     string            `json:"displayName" yaml:"displayName"`
	Jobs            []LoadTestJob     `json:"jobs" yaml:"jobs"`
	Report          LoadTestJobOutput `json:"report" yaml:"report"`
	WaitBetweenJobs int               `json:"waitBetweenJobs" yaml:"waitBetweenJobs"`
	JobType         string            `json:"type" yaml:"type"`
}
