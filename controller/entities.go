package controller

import "github.com/cjlapao/http-loadtester-go/jobs"

type LoadTestResponse struct {
	ID            string
	Name          *string
	Type          jobs.JobOPerationType
	OperationType jobs.BlockType
	Target        *jobs.JobOperationTarget
	Options       *jobs.JobOperationOptions
	Result        *jobs.JobOperationResult
}
