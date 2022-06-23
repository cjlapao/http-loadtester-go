package jobs

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/xid"
)

// JobOperationBlock Entity
type JobOperationBlock struct {
	ID              string
	JobID           string
	JobName         *string
	BlockPosition   int
	TotalBlocks     int
	Type            JobOPerationType
	JobBlockType    BlockType
	BlockType       BlockType
	Target          *JobOperationTarget
	Tasks           *[]*JobOperationBlockTask
	Result          *JobOperationBlockResult
	WaitFor         Interval
	MaxTaskInterval Interval
	MinTaskInterval Interval
	Timeout         int
}

// CreateBlock Create a Block for the JobOperation
func (j *JobOperation) CreateBlock() *JobOperationBlock {
	block := JobOperationBlock{
		ID:           xid.New().String(),
		JobID:        j.ID,
		Type:         j.Type,
		JobBlockType: j.OperationType,
		BlockType:    j.Options.BlockType,
		Target:       j.Target,
		Timeout:      j.Options.Timeout,
	}

	block.Result = block.CreateBlockResult()
	block.Result.BlockID = block.ID
	if *j.Name != "" {
		block.JobName = j.Name
	} else {
		block.JobName = &j.ID
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
			task.Execute(&taskWaitGroup)
		}
	}
	taskWaitGroup.Wait()
	endingTime := time.Now().UTC()

	for _, task := range *j.Tasks {
		*j.Result.TaskResults = append(*j.Result.TaskResults, task.Result)
	}
	var duration time.Duration = endingTime.Sub(startingTime)

	j.Result.TotalDurationInSeconds = duration.Seconds()

	logger.Info("Finished processing %v Block %v [%v/%v], made %v %v calls to target", fmt.Sprint(j.JobBlockType), j.ID, fmt.Sprint(j.BlockPosition), fmt.Sprint(j.TotalBlocks), fmt.Sprint(amountTasks), fmt.Sprint(j.BlockType))
	wg.Done()
}
