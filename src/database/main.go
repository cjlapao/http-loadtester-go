package database

import (
	"github.com/cjlapao/common-go/database/mongodb"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/http-loadtester-go/constants"
)

var MongoDbSvc *mongodb.MongoDBService
var contextSvc = execution_context.Get()

// Initializes the backend database to be used
func Init() bool {
	connectionString := contextSvc.Configuration.GetString(constants.MONGO_DB_CONNECTION_STRING)
	if connectionString == "" {
		contextSvc.Configuration.UpsertKey(constants.DATABASE_ENABLED, false)
		return false
	}

	MongoDbSvc = mongodb.Init()
	MongoDbSvc.ConnectionString = connectionString
	MongoDbSvc.GlobalDatabaseName = constants.GLOBAL_DATABASE_NAME
	contextSvc.Configuration.UpsertKey(constants.DATABASE_ENABLED, true)
	return true
}

func IsEnabled() bool {
	return contextSvc.Configuration.GetBool(constants.DATABASE_ENABLED)
}
