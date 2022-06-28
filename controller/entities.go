package controller

import "github.com/cjlapao/http-loadtester-go/domain"

type LoadTestResponse struct {
	ID            string
	Name          *string
	Type          domain.JobOPerationType
	OperationType domain.BlockType
	Target        *domain.JobOperationTarget
	Options       *domain.JobOperationOptions
	Result        *domain.JobOperationResult
}
