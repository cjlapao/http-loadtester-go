package domain

import (
	"github.com/cjlapao/common-go/configuration"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/log"
)

type JobOperationRepository interface {
	Upsert(jobOperation JobOperation)
}

type ServiceProvider interface {
	Logger() *log.Logger
	Context() *execution_context.Context
	Configuration() *configuration.ConfigurationService
	JobOperationRepo() JobOperationRepository
	IsDatabaseEnabled() bool
}

type JobOperationService interface {
	Execute(specs *JobOperationOptions) error
}
