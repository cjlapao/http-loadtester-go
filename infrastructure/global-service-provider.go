package infrastructure

import (
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/http-loadtester-go/adapters"
	"github.com/cjlapao/http-loadtester-go/domain"
)

var globalServiceProvider *GlobalServicesProvider

type GlobalServicesProvider struct{}

func GetServiceProvider() GlobalServicesProvider {
	if globalServiceProvider == nil {
		globalServiceProvider = &GlobalServicesProvider{}
	}

	return *globalServiceProvider
}

func NewServiceProvider() GlobalServicesProvider {
	globalServiceProvider = &GlobalServicesProvider{}

	return *globalServiceProvider
}

func (gp GlobalServicesProvider) Logger() *log.Logger {
	return logger
}

func (gp GlobalServicesProvider) JobOperationRepo() domain.JobOperationRepository {
	return adapters.NewJobOperationRepo()
}
