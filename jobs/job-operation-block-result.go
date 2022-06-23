package jobs

// JobOperationBlockResult Entity
type JobOperationBlockResult struct {
	JobID                  string
	BlockID                string
	TaskResults            *[]*JobOperationBlockTaskResult
	TargetCalls            map[string]int
	Total                  int
	Failed                 int
	Succeeded              int
	TaskResponseStatus     *[]*JobOperationTaskResponseStatusResult
	TotalDurationInSeconds float64
	AverageTaskDuration    float64
	ResponseDetails        *ResponseDetails
}

// Process Processes a JobOperationBlockResult updating the job
func (r *JobOperationBlockResult) Process() {
	totalDurationForAverage := 0.0
	responseStatusResult := make([]*JobOperationTaskResponseStatusResult, 0)
	r.TargetCalls = make(map[string]int)
	if r.TaskResults != nil {
		for _, callResult := range *r.TaskResults {
			r.Total++
			if val, ok := r.TargetCalls[callResult.TargetedUri]; ok {
				r.TargetCalls[callResult.TargetedUri] = val + 1
			} else {
				r.TargetCalls[callResult.TargetedUri] = 1
			}

			if callResult.StatusCode >= 200 && callResult.StatusCode <= 299 {
				r.Succeeded++
			} else {
				r.Failed++
			}
			if len(responseStatusResult) == 0 {
				newStatusCode := JobOperationTaskResponseStatusResult{
					Code:  callResult.StatusCode,
					Count: 1,
				}
				responseStatusResult = append(responseStatusResult, &newStatusCode)
			} else {
				exists := false
				for _, status := range responseStatusResult {
					if status.Code == callResult.StatusCode {
						status.Count++
						exists = true
						break
					}
				}

				if !exists {
					newStatusCode := JobOperationTaskResponseStatusResult{
						Code:  callResult.StatusCode,
						Count: 1,
					}
					responseStatusResult = append(responseStatusResult, &newStatusCode)
				}
			}

			totalDurationForAverage += callResult.QueryDuration.Seconds
			if r.ResponseDetails == nil && callResult.ResponseDetails != nil {
				r.ResponseDetails = callResult.ResponseDetails
			}
		}

		r.TaskResponseStatus = &responseStatusResult
		r.AverageTaskDuration = totalDurationForAverage / float64(r.Total)
	}
}
