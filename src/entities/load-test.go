package entities

type LoadTest struct {
	DisplayName     string    `json:"displayName" yaml:"displayName"`
	Jobs            []Job     `json:"jobs" yaml:"jobs"`
	Report          JobOutput `json:"report" yaml:"report"`
	WaitBetweenJobs int       `json:"waitBetweenJobs" yaml:"waitBetweenJobs"`
	JobType         string    `json:"type" yaml:"type"`
}
