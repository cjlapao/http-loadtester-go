package jobs

import "github.com/cjlapao/http-loadtester-go/common"

// JobOperationTarget Definition
type JobOperationTarget struct {
	URLs                 []string
	Method               TargetMethod
	ContentType          string
	Body                 string
	JwtTokens            []string
	BasicAuthentications []string
	UserAgent            string
	Headers              map[string]string
	logResponse          bool
}

// CreateTarget Creates a Default Target to the JobOperation
func (j *JobOperation) CreateTarget() *JobOperationTarget {
	target := JobOperationTarget{
		Method:               GET,
		UserAgent:            "http-load-tester",
		Headers:              make(map[string]string),
		URLs:                 make([]string, 0),
		JwtTokens:            make([]string, 0),
		BasicAuthentications: make([]string, 0),
	}

	j.Target = &target
	return j.Target
}

// SetJSONContent Sets the operation type to JSON
func (t *JobOperationTarget) SetJSONContent() {
	t.ContentType = "application/json; charset=UTF-8"
}

func (j *JobOperationTarget) IsMultiTargeted() bool {
	if j.URLs != nil && len(j.URLs) > 1 {
		return true
	}

	return false
}

func (j *JobOperationTarget) GetUrl(index int) string {
	if len(j.URLs) < index {
		return ""
	}

	return j.URLs[index]
}

func (j *JobOperationTarget) CountUrls() int {
	return len(j.URLs)
}

func (j *JobOperationTarget) GetRandomUrl() string {
	if j.URLs != nil && len(j.URLs) == 1 {
		return j.URLs[0]
	}

	randIndex := common.GetRandomNum(0, len(j.URLs))

	return j.URLs[randIndex]
}

func (j *JobOperationTarget) HasJwtAuthentication() bool {
	if j.JwtTokens != nil && len(j.JwtTokens) > 0 {
		return true
	}

	return false
}

func (j *JobOperationTarget) GetRandomJwtToken() string {
	if j.JwtTokens != nil && len(j.JwtTokens) == 1 {
		return j.JwtTokens[0]
	}

	randIndex := common.GetRandomNum(0, len(j.JwtTokens))

	return j.JwtTokens[randIndex]
}

func (j *JobOperationTarget) HasBasicAuthentication() bool {
	if j.BasicAuthentications != nil && len(j.BasicAuthentications) > 0 {
		return true
	}

	return false
}

func (j *JobOperationTarget) GetRandomBasicAuthentication() string {
	if j.BasicAuthentications != nil && len(j.BasicAuthentications) == 1 {
		return j.BasicAuthentications[0]
	}

	randIndex := common.GetRandomNum(0, len(j.BasicAuthentications))

	return j.BasicAuthentications[randIndex]
}

func (j *JobOperationTarget) HasHeaders() bool {
	if j.Headers != nil && len(j.Headers) > 0 {
		return true
	}

	return false
}
