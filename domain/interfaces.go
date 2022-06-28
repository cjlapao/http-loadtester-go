package domain

import "github.com/cjlapao/common-go/log"

type JobOperationRepository interface {
	Upsert(jobOperation JobOperation)
}

type ServiceProvider interface {
	Logger() *log.Logger
	JobOperationRepo() JobOperationRepository
}
