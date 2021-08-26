package jobs

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

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

// JobOperationOptions Entity
type JobOperationOptions struct {
	Duration         int
	MaxTaskOutput    int
	Timeout          int
	Verbose          bool
	LogResult        bool
	BlockType        BlockType
	BlockInterval    Interval
	TasksPerBlock    Interval
	MaxBlockInterval Interval
	MinBlockInterval Interval
	MaxTaskInterval  Interval
	MinTaskInterval  Interval
	MaxTasksPerBlock Interval
	MinTasksPerBlock Interval
}

// JobOperationResult Entity
type JobOperationResult struct {
	ID                     string
	Target                 *JobOperationTarget
	BlockResults           *[]*JobOperationBlockResult
	MaxTaskOutput          int
	Total                  int
	TotalCalls             int
	TotalSucceededCalls    int
	TotalFailedCalls       int
	TotalDurationInSeconds float64
	AverageBlockDuration   float64
	AverageCallDuration    float64
	ResponseDetails        *ResponseDetails
	TimeTaken              time.Duration
	TaskResponseStatus     *[]*JobOperationTaskResponseStatusResult
}

// CreateJobOperation Creates a new Job Operation Task
func CreateJobOperation() *JobOperation {
	job := JobOperation{
		ID:            xid.New().String(),
		OperationType: ParallelBlock,
		Options: &JobOperationOptions{
			Timeout:          120,
			Verbose:          false,
			Duration:         60,
			MaxTaskOutput:    15,
			MaxBlockInterval: NewInterval(0),
			MinBlockInterval: NewInterval(0),
			MaxTasksPerBlock: NewInterval(20),
			MinTasksPerBlock: NewInterval(10),
			MaxTaskInterval:  NewInterval(0),
			MinTaskInterval:  NewInterval(0),
			BlockInterval:    NewInterval(1),
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
	numberOfBlocks := 0
	switch j.Type {
	case Constant, Increasing:
		if j.Options.BlockInterval.Value() <= 0 {
			j.Options.BlockInterval = NewInterval(1)
		}
		numberOfBlocks = j.Options.Duration / j.Options.BlockInterval.Value()
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
		numberOfBlocks = j.Options.Duration / j.Options.MaxBlockInterval.Value()
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

	randomBlockNumber := callRandom.Intn(max-min) + min

	return randomBlockNumber
}

func (j *JobOperation) getRandomTaskCount() int {
	max := j.Options.MaxTasksPerBlock.Value()
	min := j.Options.MinTasksPerBlock.value
	rand.Seed(time.Now().UnixNano())
	someSalt := int64(rand.Intn(10000))
	rand.Seed(time.Now().UnixNano() * someSalt)
	randomTasksNumber := rand.Intn(max-min) + min

	return randomTasksNumber
}

// Execute Executes a Job Operation creating X amount of blocks that will be run every X seconds
// This will be defined by the amount of blocks the duration of the load test
func (j *JobOperation) Execute() error {
	j.generateBlocks()
	amountOfBlocks := len(j.Blocks)
	j.Result = j.CreateLoadJobResult()
	var blockWaitingGroup sync.WaitGroup
	blockWaitingGroup.Add(amountOfBlocks)
	logger.Command("Performing a Load Test on %v for %v seconds\nThis can take longer depending on the pressure of the tasks being performed", j.Target.URL, fmt.Sprint(j.Options.Duration))
	startingTime := time.Now().UTC()
	// Executing the blocks
	for i, block := range j.Blocks {
		blockNum := i + 1
		block.BlockPosition = blockNum
		block.TotalBlocks = amountOfBlocks
		switch j.OperationType {
		case ParallelBlock:
			logger.Info("Started processing Parallel Block %v [%v/%v], using %v load with %v %v tasks and %v timeout", block.ID, fmt.Sprint(blockNum), fmt.Sprint(amountOfBlocks), fmt.Sprint(j.Type), fmt.Sprint(len(*block.Tasks)), fmt.Sprint(block.BlockType), fmt.Sprint(time.Duration(j.Options.Timeout)*time.Second))
			go block.Execute(&blockWaitingGroup)
		case SequentialBlock:
			logger.Info("Started processing Sequential Block %v [%v/%v], using %v load with %v %v tasks and %v timeout", block.ID, fmt.Sprint(blockNum), fmt.Sprint(amountOfBlocks), fmt.Sprint(j.Type), fmt.Sprint(len(*block.Tasks)), fmt.Sprint(block.BlockType), fmt.Sprint(time.Duration(j.Options.Timeout)*time.Second))
			block.Execute(&blockWaitingGroup)
		default:
			logger.Info("Started processing Sequential Block %v [%v/%v], using %v load with %v %v tasks and %v timeout", block.ID, fmt.Sprint(blockNum), fmt.Sprint(amountOfBlocks), fmt.Sprint(j.Type), fmt.Sprint(len(*block.Tasks)), fmt.Sprint(block.BlockType), fmt.Sprint(time.Duration(j.Options.Timeout)*time.Second))
			block.Execute(&blockWaitingGroup)
		}
		if block.WaitFor.Value() > 0 && i < len(j.Blocks) {
			logger.Info("Waiting for %v before starting next block", fmt.Sprint(time.Duration(block.WaitFor.Value())*time.Second))
			time.Sleep(time.Duration(block.WaitFor.Value()) * time.Second)
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
	logger.Success("Finished Load Test on %v for %v seconds", j.Target.URL, fmt.Sprint(j.Options.Duration))
	return nil
}

// Authenticated Check if this job is using Authentication
func (j *JobOperation) Authenticated() bool {
	if j.Target == nil {
		return false
	}
	if j.Target.JwtToken != "" {
		return true
	}
	return false
}

// ProcessResult Processes the results and calculates the averages
func (j *JobOperationResult) ProcessResult() {
	totalDurationForAverage := 0.0
	totalTasksDurationForAverage := 0.0
	responseStatusResult := make([]*JobOperationTaskResponseStatusResult, 0)
	if j.BlockResults != nil {
		for _, blockResult := range *j.BlockResults {
			j.Total++
			j.TotalCalls += blockResult.Total
			j.TotalSucceededCalls += blockResult.Succeeded
			j.TotalFailedCalls += blockResult.Failed
			totalDurationForAverage += blockResult.TotalDurationInSeconds

			for _, status := range *blockResult.TaskResponseStatus {
				exists := false
				for _, jobStatus := range responseStatusResult {
					if status.Code == jobStatus.Code {
						jobStatus.Count = jobStatus.Count + status.Count
						exists = true
					}
				}

				if !exists {
					newStatusCode := JobOperationTaskResponseStatusResult{
						Code:  status.Code,
						Count: status.Count,
					}

					responseStatusResult = append(responseStatusResult, &newStatusCode)
				}
			}

			for _, taskResult := range *blockResult.TaskResults {
				totalTasksDurationForAverage += taskResult.QueryDuration.Seconds
				if j.ResponseDetails == nil && blockResult.ResponseDetails != nil {
					j.ResponseDetails = blockResult.ResponseDetails
				}
			}

		}

		j.TaskResponseStatus = &responseStatusResult
		j.AverageBlockDuration = totalDurationForAverage / float64(j.Total)
		j.AverageCallDuration = totalTasksDurationForAverage / float64(j.TotalCalls)
	}
}
