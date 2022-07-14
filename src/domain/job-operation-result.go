package domain

import "time"

// JobOperationResult Entity
type JobOperationResult struct {
	ID                       string
	Target                   *JobOperationTarget
	BlockResults             *[]*JobOperationBlockResult
	TargetCalls              map[string]int
	TargetAuthenticationUsed map[string]map[string]int
	MaxTaskOutput            int
	Total                    int
	TotalCalls               int
	TotalSucceededCalls      int
	TotalFailedCalls         int
	TotalDurationInSeconds   float64
	AverageBlockDuration     float64
	AverageCallDuration      float64
	ResponseDetails          *ResponseDetails
	StartingTime             time.Time
	EndingTime               time.Time
	TimeTaken                time.Duration
	TaskResponseStatus       *[]*JobOperationTaskResponseStatusResult
}

// ProcessResult Processes the results and calculates the averages
func (j *JobOperationResult) ProcessResult() {
	totalDurationForAverage := 0.0
	totalTasksDurationForAverage := 0.0
	j.TargetCalls = make(map[string]int)
	j.TargetAuthenticationUsed = make(map[string]map[string]int)
	responseStatusResult := make([]*JobOperationTaskResponseStatusResult, 0)
	if j.BlockResults != nil {
		for _, blockResult := range *j.BlockResults {
			j.Total++
			j.TotalCalls += blockResult.Total
			j.TotalSucceededCalls += blockResult.Succeeded
			j.TotalFailedCalls += blockResult.Failed
			totalDurationForAverage += blockResult.TotalDurationInSeconds

			for key, blockVal := range blockResult.TargetCalls {
				if val, ok := j.TargetCalls[key]; ok {
					j.TargetCalls[key] = blockVal + val
				} else {
					j.TargetCalls[key] = blockVal
				}
			}

			// recording the usage of the authentication
			for key, blockVal := range blockResult.TargetAuthenticationUsed {
				if _, ok := j.TargetAuthenticationUsed[key]; ok {
					for blockTargetKey, blockTargetVal := range blockVal {
						if _, ok := j.TargetAuthenticationUsed[key][blockTargetKey]; ok {
							j.TargetAuthenticationUsed[key][blockTargetKey] = j.TargetAuthenticationUsed[key][blockTargetKey] + blockTargetVal
						} else {
							j.TargetAuthenticationUsed[key][blockTargetKey] = blockTargetVal
						}
					}
				} else {
					j.TargetAuthenticationUsed[key] = blockVal
				}
			}

			for _, status := range *blockResult.TaskResponseStatus {
				exists := false
				for _, jobStatus := range responseStatusResult {
					if status.Code == jobStatus.Code {
						jobStatus.Count = jobStatus.Count + status.Count
						exists = true
					}
				}

				if !exists {
					newStatusCode := JobOperationTaskResponseStatusResult{
						Code:  status.Code,
						Count: status.Count,
					}

					responseStatusResult = append(responseStatusResult, &newStatusCode)
				}
			}

			for _, taskResult := range *blockResult.TaskResults {
				totalTasksDurationForAverage += taskResult.QueryDuration.Seconds
				if j.ResponseDetails == nil && blockResult.ResponseDetails != nil {
					j.ResponseDetails = blockResult.ResponseDetails
				}
			}

		}

		j.TaskResponseStatus = &responseStatusResult
		j.AverageBlockDuration = totalDurationForAverage / float64(j.Total)
		j.AverageCallDuration = totalTasksDurationForAverage / float64(j.TotalCalls)
	}
}
