package jobs

import (
	"fmt"
	"sync"
	"time"

	"github.com/cjlapao/http-loadtester-go/common"
	"github.com/rs/xid"
)

// JobOperation Constants
const (
	IncreaseFactor float64 = 0.03
)

// JobOperation Entity
type JobOperation struct {
	ID            string
	Name          *string
	Type          JobOPerationType
	OperationType BlockType
	Target        *JobOperationTarget
	Options       *JobOperationOptions
	Result        *JobOperationResult
	Blocks        []*JobOperationBlock
}

// CreateJobOperation Creates a new Job Operation Task
func CreateJobOperation() *JobOperation {
	job := JobOperation{
		ID:            xid.New().String(),
		OperationType: ParallelBlock,
		Options: &JobOperationOptions{
			Timeout:          120000,
			Verbose:          false,
			Duration:         60,
			MaxTaskOutput:    15,
			MaxBlockInterval: NewInterval(0),
			MinBlockInterval: NewInterval(0),
			MaxTasksPerBlock: NewInterval(20),
			MinTasksPerBlock: NewInterval(10),
			MaxTaskInterval:  NewInterval(0),
			MinTaskInterval:  NewInterval(0),
			BlockInterval:    NewInterval(0),
			TasksPerBlock:    NewInterval(60),
		},
	}

	job.Name = &job.ID
	job.CreateTarget()
	job.CreateLoadJobResult()
	return &job
}

// CreateLoadJobResult Creates a JobOperation result
func (j *JobOperation) CreateLoadJobResult() *JobOperationResult {
	jobResults := make([]*JobOperationBlockResult, 0)
	result := JobOperationResult{
		ID:            j.ID,
		Target:        j.Target,
		MaxTaskOutput: j.Options.MaxTaskOutput,
		BlockResults:  &jobResults,
	}

	return &result
}

func (j *JobOperation) generateBlocks() {
	numberOfBlocks := j.Options.NumberOfBlocks

	if numberOfBlocks <= 0 {
		numberOfBlocks = 1
	}

	switch j.Type {
	case Constant, Increasing:
		if j.Options.BlockInterval.Value() <= 0 {
			j.Options.BlockInterval = NewInterval(0)
		}
		if numberOfBlocks > 0 {
			for i := 0; i < numberOfBlocks; i++ {
				block := j.CreateBlock()
				block.WaitFor = NewInterval(j.Options.BlockInterval.Value())
				block.BlockType = j.Options.BlockType
				block.MinTaskInterval = j.Options.MinTaskInterval
				block.MaxTaskInterval = j.Options.MaxTaskInterval
			}
		}
	case Fuzzy:
		if j.Options.MaxBlockInterval.Value() <= 0 {
			j.Options.MaxBlockInterval = NewInterval(1)
			j.Options.MinBlockInterval = NewInterval(1)
		}
		if numberOfBlocks > 0 {
			for i := 0; i < numberOfBlocks; i++ {
				block := j.CreateBlock()
				block.BlockType = j.Options.BlockType
				block.MinTaskInterval = j.Options.MinTaskInterval
				block.MaxTaskInterval = j.Options.MaxTaskInterval
				if j.Options.MaxBlockInterval.Value() > j.Options.MinBlockInterval.value {
					block.WaitFor = NewInterval(j.getRandomBlockInterval())
				} else {
					block.WaitFor = NewInterval(j.Options.MaxBlockInterval.Value())
				}
			}
		}
	}
	j.generateBlockTasks()
}

func (j *JobOperation) generateBlockTasks() {
	previousTaskCount := 0.0
	blockSize := float64(len(j.Blocks))
	numberTasks := float64(j.Options.TasksPerBlock.Value())
	numberTasksPerBlock := numberTasks / blockSize
	factor := numberTasks / blockSize / float64(j.Options.Duration)
	tasksWithFactor := numberTasksPerBlock * factor
	for _, block := range j.Blocks {
		switch j.Type {
		case Constant:
			if j.Options.TasksPerBlock.Value() <= 0 {
				j.Options.TasksPerBlock = NewInterval(10)
			}
			tasksPerBlock := j.Options.TasksPerBlock.Value()
			for nTask := 0; nTask < tasksPerBlock; nTask++ {
				block.CreateTask(nTask + 1)
			}
		case Increasing:
			if j.Options.TasksPerBlock.Value() <= 0 {
				j.Options.TasksPerBlock = NewInterval(10)
			}
			tasksPerBlock := int(tasksWithFactor + previousTaskCount)
			if tasksPerBlock <= 0 {
				tasksPerBlock = 1
			}
			previousTaskCount = float64(tasksPerBlock)
			for nTask := 0; nTask < tasksPerBlock; nTask++ {
				block.CreateTask(nTask + 1)
			}
		case Fuzzy:
			if j.Options.MaxTasksPerBlock.Value() <= 0 {
				j.Options.MaxTasksPerBlock = NewInterval(10)
				j.Options.MinTasksPerBlock = NewInterval(1)
			}
			tasksPerBlock := j.Options.MaxTasksPerBlock.Value()
			if j.Options.MaxTasksPerBlock.Value() > j.Options.MinTasksPerBlock.value {
				tasksPerBlock = j.getRandomTaskCount()
			}
			for nTask := 0; nTask < tasksPerBlock; nTask++ {
				block.CreateTask(nTask + 1)
			}
		}

	}
}

func (j *JobOperation) getRandomBlockInterval() int {
	max := j.Options.MaxBlockInterval.Value()
	min := j.Options.MinBlockInterval.Value()

	return common.GetRandomNum(min, max)
}

func (j *JobOperation) getRandomTaskCount() int {
	max := j.Options.MaxTasksPerBlock.Value()
	min := j.Options.MinTasksPerBlock.value

	return common.GetRandomNum(min, max)
}

// Execute Executes a Job Operation creating X amount of blocks that will be run every X seconds
// This will be defined by the amount of blocks the duration of the load test
func (j *JobOperation) Execute(wg *sync.WaitGroup) error {
	startingJobTime := time.Now().UTC()
	j.generateBlocks()
	amountOfBlocks := len(j.Blocks)
	j.Result = j.CreateLoadJobResult()
	var blockWaitingGroup sync.WaitGroup
	blockWaitingGroup.Add(amountOfBlocks)
	if j.Target.IsMultiTargeted() {
		logger.Success("Performing Load Test on %v targets with %v blocks", fmt.Sprintf("%v", j.Target.CountUrls()), fmt.Sprint(j.Options.NumberOfBlocks))
	} else {
		logger.Success("Performing Load Test on %v with %v blocks", j.Target.GetUrl(0), fmt.Sprint(j.Options.NumberOfBlocks))
	}

	startingTime := time.Now().UTC()
	// Executing the blocks
	for i, block := range j.Blocks {
		blockNum := i + 1
		block.BlockPosition = blockNum
		block.TotalBlocks = amountOfBlocks
		logger.Info("Started processing %v Block %v [%v/%v], using %v load with %v %v tasks and %v timeout", fmt.Sprint(j.OperationType), block.ID, fmt.Sprint(blockNum), fmt.Sprint(amountOfBlocks), fmt.Sprint(j.Type), fmt.Sprint(len(*block.Tasks)), fmt.Sprint(block.BlockType), fmt.Sprint(time.Duration(j.Options.Timeout)*time.Millisecond))
		switch j.OperationType {
		case ParallelBlock:
			go block.Execute(&blockWaitingGroup)
		case SequentialBlock:
			block.Execute(&blockWaitingGroup)
		default:
			block.Execute(&blockWaitingGroup)
		}
		if block.WaitFor.Value() > 0 && i < len(j.Blocks) {
			logger.Info("Waiting for %v before starting next block", fmt.Sprint(time.Duration(block.WaitFor.Value())*time.Millisecond))
			time.Sleep(time.Duration(block.WaitFor.Value()) * time.Millisecond)
		}
	}

	blockWaitingGroup.Wait()
	endingTime := time.Now().UTC()
	var duration time.Duration = endingTime.Sub(startingTime)

	j.Result.TotalDurationInSeconds = duration.Seconds()
	// Parsing the results after we go all of them done
	for _, block := range j.Blocks {
		block.Result.Process()
		*j.Result.BlockResults = append(*j.Result.BlockResults, block.Result)
	}

	j.Result.ProcessResult()
	if j.Target.IsMultiTargeted() {
		logger.Success("Finished Load Test on %v targets for %v seconds", fmt.Sprintf("%v", j.Target.CountUrls()), fmt.Sprint(j.Options.Duration))

	} else {
		logger.Success("Finished Load Test on %v for %v seconds", j.Target.GetUrl(0), fmt.Sprint(j.Options.Duration))
	}
	endingJobTime := time.Now().UTC()
	j.Result.TimeTaken = endingJobTime.Sub(startingJobTime)
	j.Result.StartingTime = startingJobTime
	j.Result.EndingTime = endingJobTime
	wg.Done()
	return nil
}

// Authenticated Check if this job is using Authentication
func (j *JobOperation) Authenticated() bool {
	if j.Target == nil {
		return false
	}
	if j.Target.JwtTokens != nil || len(j.Target.JwtTokens) > 0 {
		return true
	}
	return false
}
