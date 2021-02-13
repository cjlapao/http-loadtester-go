package jobrunner

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/xid"
)

// JobOperationBlock Entity
type JobOperationBlock struct {
	ID        string
	JobID     string
	JobName   string
	Type      JobOPerationType
	BlockType BlockType
	Target    *JobOperationTarget
	Tasks     *[]*JobOperationBlockTask
	Result    *JobOperationBlockResult
	WaitFor   Interval
	Timeout   int
}

// CreateBlock Create a Block for the JobOperation
func (j *JobOperation) CreateBlock() *JobOperationBlock {
	block := JobOperationBlock{
		ID:        xid.New().String(),
		JobID:     j.ID,
		Type:      j.Type,
		BlockType: j.BlockType,
		Target:    j.Target,
		Timeout:   j.Options.Timeout,
	}
	block.Result = block.CreateBlockResult()
	block.Result.BlockID = block.ID
	if j.Name != "" {
		block.JobName = j.Name
	} else {
		block.JobName = j.ID
	}
	if j.Blocks == nil {
		j.Blocks = make([]*JobOperationBlock, 0)
	}

	j.Blocks = append(j.Blocks, &block)
	return &block
}

// CreateBlockResult Create a JobOperationBlock result
func (j *JobOperationBlock) CreateBlockResult() *JobOperationBlockResult {
	taskResults := make([]*JobOperationBlockTaskResult, 0)
	result := JobOperationBlockResult{
		TaskResults: &taskResults,
	}
	result.BlockID = j.ID
	result.JobID = j.JobID
	return &result
}

// Execute Executes a Block
func (j *JobOperationBlock) Execute(wg *sync.WaitGroup) {
	amountTasks := len(*j.Tasks)
	logger.Debug("Executing %v tasks for block %v", fmt.Sprint(amountTasks), j.ID)
	var taskWaitGroup sync.WaitGroup
	taskWaitGroup.Add(amountTasks)
	startingTime := time.Now().UTC()

	for _, task := range *j.Tasks {
		switch j.BlockType {
		case ParallelBlock:
			go task.Execute(&taskWaitGroup)
		case SequentialBlock:
			task.Execute(&taskWaitGroup)
		default:
			go task.Execute(&taskWaitGroup)
		}
	}
	taskWaitGroup.Wait()
	endingTime := time.Now().UTC()

	for _, task := range *j.Tasks {
		*j.Result.TaskResults = append(*j.Result.TaskResults, task.Result)
	}
	var duration time.Duration = endingTime.Sub(startingTime)

	j.Result.TotalDurationInSeconds = duration.Seconds()

	logger.Info("Finished processing Block %v, made %v calls to target", j.ID, fmt.Sprint(amountTasks))
	wg.Done()
}

// JobOperationBlockResult Entity
type JobOperationBlockResult struct {
	JobID                  string
	BlockID                string
	TaskResults            *[]*JobOperationBlockTaskResult
	Total                  int
	Failed                 int
	Succeeded              int
	TotalDurationInSeconds float64
	AverageTaskDuration    float64
}

// Process Processes a JobOperationBlockResult updating the job
func (r *JobOperationBlockResult) Process() {
	totalDurationForAverage := 0.0
	if r.TaskResults != nil {
		for _, callResult := range *r.TaskResults {
			r.Total++
			if callResult.StatusCode >= 200 && callResult.StatusCode <= 299 {
				r.Succeeded++
			} else {
				r.Failed++
			}
			totalDurationForAverage += callResult.QueryDuration.Seconds
		}
		r.AverageTaskDuration = totalDurationForAverage / float64(r.Total)
	}
}
