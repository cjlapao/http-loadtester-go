package adapters

import (
	"github.com/cjlapao/common-go/database/mongodb"
	"github.com/cjlapao/http-loadtester-go/constants"
	"github.com/cjlapao/http-loadtester-go/database"
	"github.com/cjlapao/http-loadtester-go/domain"
)

type JobOperationRepo struct {
	Repo mongodb.MongoRepository
}

func NewJobOperationRepo() JobOperationRepo {
	result := JobOperationRepo{}
	result.Repo = database.MongoDbSvc.GlobalDatabase().NewRepository(constants.JOB_OPERATION_REPOSITORY_NAME)
	return result
}

func (repo JobOperationRepo) Upsert(jobOperation domain.JobOperation) {
	model, _ := mongodb.NewUpdateOneModelBuilder().Encode(jobOperation).FilterBy("id", mongodb.Equal, jobOperation.ID).Build()

	repo.Repo.UpsertOne(model)
}
