package jobrunner

// JobOperationTarget Definition
type JobOperationTarget struct {
	URL         string
	Method      TargetMethod
	ContentType string
	Body        string
	JwtToken    string
}

// CreateTarget Creates a Default Target to the JobOperation
func (j *JobOperation) CreateTarget() *JobOperationTarget {
	target := JobOperationTarget{
		Method: GET,
	}

	j.Target = &target
	return j.Target
}

func (t *JobOperationTarget) SetJsonContent() {
	t.ContentType = "application/json; charset=UTF-8"
}
