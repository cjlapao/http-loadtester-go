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
		contextSvc.Configuration.UpsertKey("isMongoEnabled", false)
		return false
	}

	MongoDbSvc = mongodb.Init()
	MongoDbSvc.ConnectionString = connectionString
	MongoDbSvc.GlobalDatabaseName = "http_load_tester"
	contextSvc.Configuration.UpsertKey("isMongoEnabled", true)
	return true
}

func IsEnabled() bool {
	return contextSvc.Configuration.GetBool("isMongoEnabled")
}
