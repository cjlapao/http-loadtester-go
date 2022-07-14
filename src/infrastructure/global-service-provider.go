package infrastructure

import (
	"github.com/cjlapao/common-go/configuration"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/http-loadtester-go/adapters"
	"github.com/cjlapao/http-loadtester-go/constants"
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

func (gp GlobalServicesProvider) Context() *execution_context.Context {
	return execution_context.Get()
}

func (gp GlobalServicesProvider) Configuration() *configuration.ConfigurationService {
	return gp.Context().Configuration
}

func (gp GlobalServicesProvider) JobOperationRepo() domain.JobOperationRepository {
	return adapters.NewJobOperationRepo()
}

func (gp GlobalServicesProvider) IsDatabaseEnabled() bool {
	return gp.Configuration().GetBool(constants.DATABASE_ENABLED)
}
