package jobs

// JobOperationTarget Definition
type JobOperationTarget struct {
	URL         string
	Method      TargetMethod
	ContentType string
	Body        string
	JwtToken    string
	logResponse bool
}

// CreateTarget Creates a Default Target to the JobOperation
func (j *JobOperation) CreateTarget() *JobOperationTarget {
	target := JobOperationTarget{
		Method: GET,
	}

	j.Target = &target
	return j.Target
}

// SetJSONContent Sets the operation type to JSON
func (t *JobOperationTarget) SetJSONContent() {
	t.ContentType = "application/json; charset=UTF-8"
}
