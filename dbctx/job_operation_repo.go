package dbctx

import (
	"github.com/cjlapao/common-go/database/mongodb"
	"github.com/cjlapao/http-loadtester-go/constants"
	"github.com/cjlapao/http-loadtester-go/jobs"
)

var repo *JobOperationRepository

type JobOperationRepository struct {
	Repository mongodb.MongoRepository
}

func Get() *JobOperationRepository {

	if repo == nil {
		repo = &JobOperationRepository{}
	}

	repo.Repository = mongodb.Get().GlobalDatabase().NewRepository(constants.JOB_OPERATION_REPOSITORY_NAME)

	return repo
}

func (repo *JobOperationRepository) Upsert(operation jobs.JobOperation) error {
	model, err := mongodb.NewUpdateOneModelBuilder().FilterBy("id", mongodb.Equal, operation.ID).Build()
	if err != nil {
		return err
	}

	repo.Repository.UpsertOne(model)
	return nil
}
