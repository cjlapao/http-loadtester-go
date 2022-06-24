package jobs

// JobOperationBlockTaskResult Entity
type JobOperationBlockTaskResult struct {
	TaskID             int
	BlockID            string
	JobID              string
	TargetedUri        string
	AuthenticationUsed string
	Target             *JobOperationTarget
	QueryDuration      *JobOperationBlockTaskDuration
	Status             string
	StatusCode         int
	Content            string
	ErrorMessage       string
	ResponseDetails    *ResponseDetails
}
