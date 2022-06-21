package jobs

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cjlapao/common-go/helper"

	"gopkg.in/yaml.v2"
)

// LoadTest Entity
type LoadTest struct {
	DisplayName     string            `json:"displayName" yaml:"displayName"`
	Jobs            []LoadTestJob     `json:"jobs" yaml:"jobs"`
	Report          LoadTestJobOutput `json:"report" yaml:"report"`
	WaitBetweenJobs int               `json:"waitBetweenJobs" yaml:"waitBetweenJobs"`
	JobType         string            `json:"type" yaml:"type"`
}

// LoadTestJob Entity
type LoadTestJob struct {
	Name           string                 `json:"name" yaml:"name"`
	Type           string                 `json:"type" yaml:"type"`
	Target         LoadTestJobTarget      `json:"target" yaml:"target"`
	ConstantLoad   *LoadTestConstantJob   `json:"constantLoad" yaml:"constantLoad"`
	IncreasingLoad *LoadTestIncreasingJob `json:"increasingLoad" yaml:"increasingLoad"`
	FuzzyLoad      *LoadTestFuzzyJob      `json:"fuzzyLoad" yaml:"fuzzyLoad"`
}

// LoadTestJobTarget Entity
type LoadTestJobTarget struct {
	URL          string   `json:"url" yaml:"url"`
	Method       string   `json:"method" yaml:"method"`
	Body         string   `json:"body" yaml:"body"`
	BearerToken  string   `json:"token" yaml:"token"`
	BearerTokens []string `json:"tokens" yaml:"tokens"`
	ContentType  string   `json:"contentType" yaml:"contentType"`
	Timeout      int      `json:"timeout" yaml:"timeout"`
	LogResponse  bool     `json:"logResponse" yaml:"logResponse"`
}

// LoadTestConstantJob entity
type LoadTestConstantJob struct {
	Duration int                        `json:"duration" yaml:"duration"`
	Options  LoadTestConstantJobOptions `json:"specs" yaml:"specs"`
}

// LoadTestConstantJobOptions Entity
type LoadTestConstantJobOptions struct {
	BlockType       string `json:"type" yaml:"type"`
	BlockInterval   int    `json:"blockInterval" yaml:"blockInterval"`
	MaxTaskInterval int    `json:"maxTaskInterval" yaml:"maxTaskInterval"`
	MinTaskInterval int    `json:"minTaskInterval" yaml:"minTaskInterval"`
	CallsPerBlock   int    `json:"callsPerBlock" yaml:"callsPerBlock"`
}

// LoadTestIncreasingJob Entity
type LoadTestIncreasingJob struct {
	Duration int                          `json:"duration" yaml:"duration"`
	Options  LoadTestIncreasingJobOptions `json:"specs" yaml:"specs"`
}

// LoadTestIncreasingJobOptions Entity
type LoadTestIncreasingJobOptions struct {
	BlockType       string `json:"type" yaml:"type"`
	BlockInterval   int    `json:"blockInterval" yaml:"blockInterval"`
	TotalCalls      int    `json:"totalCalls" yaml:"totalCalls"`
	MaxTaskInterval int    `json:"maxTaskInterval" yaml:"maxTaskInterval"`
	MinTaskInterval int    `json:"minTaskInterval" yaml:"minTaskInterval"`
}

// LoadTestFuzzyJob Entity
type LoadTestFuzzyJob struct {
	Duration int                     `json:"duration" yaml:"duration"`
	Options  LoadTestFuzzyJobOptions `json:"specs" yaml:"specs"`
}

// LoadTestFuzzyJobOptions Entity
type LoadTestFuzzyJobOptions struct {
	BlockType        string `json:"type" yaml:"type"`
	MaxBlockInterval int    `json:"maxBlockInterval" yaml:"maxBlockInterval"`
	MinBlockInterval int    `json:"minBlockInterval" yaml:"minBlockInterval"`
	MaxTasksPerBlock int    `json:"maxTasksPerBlock" yaml:"maxTasksPerBlock"`
	MinTasksPerBlock int    `json:"minTasksPerBlock" yaml:"minTasksPerBlock"`
	MaxTaskInterval  int    `json:"maxTaskInterval" yaml:"maxTaskInterval"`
	MinTaskInterval  int    `json:"minTaskInterval" yaml:"minTaskInterval"`
}

// LoadTestJobOutput Entity
type LoadTestJobOutput struct {
	MaxTaskOutput    int    `json:"maxTaskOutput" yaml:"maxTaskOutput"`
	OutputResults    bool   `json:"outputResults" yaml:"outputResults"`
	OutputToFile     bool   `json:"outputToFile" yaml:"outputToFile"`
	OutputToFilePath string `json:"outputToFilePath" yaml:"outputToFilePath"`
}

// ExecuteFromFile Execute LoadTest from file
func ExecuteFromFile(filepath string) error {
	if !helper.FileExists(filepath) {
		err := errors.New("File was not found")
		logger.Error(err.Error())
		return err
	}

	content, err := helper.ReadFromFile(filepath)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var loadTest LoadTest
	err = yaml.Unmarshal(content, &loadTest)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	jobs, err := ExecuteLoadTest(loadTest)
	if err != nil {
		return err
	}

	if loadTest.Report.OutputToFile {
		for _, job := range jobs {
			job.ExportReportToFile(loadTest.Report.OutputToFilePath)
			logger.Success("Finished creating reports for job %v", *job.Name)
		}
	} else {
		for _, job := range jobs {
			fmt.Println(job.MarkDown())
			logger.Success("Finished creating reports for job %v", *job.Name)
		}
	}
	if loadTest.Report.OutputResults {
		for _, job := range jobs {

			job.ExportOutputToFile(loadTest.Report.OutputToFilePath)
			logger.Success("Finished creating reports for job %v", *job.Name)
		}
	}

	return nil
}

func ExecuteLoadTest(loadTest LoadTest) ([]*JobOperation, error) {
	if loadTest.DisplayName != "" {
		logger.Success("Testing File %v ", loadTest.DisplayName)
	}
	loadTesterJobs := make([]*JobOperation, 0)

	for _, loadTesterJob := range loadTest.Jobs {
		if loadTesterJob.Target.URL == "" {
			err := errors.New("Url was not defined")
			logger.Error(err.Error())
			return loadTesterJobs, err
		}

		job := CreateJobOperation()

		switch strings.ToLower(loadTesterJob.Type) {
		case "parallel":
			job.OperationType = ParallelBlock
		case "sequential":
			job.OperationType = SequentialBlock
		default:
			job.OperationType = SequentialBlock
		}
		if loadTesterJob.Name != "" {
			jName := strings.ReplaceAll(loadTesterJob.Name, " ", "_")
			jName = strings.ReplaceAll(jName, ":", "")
			jName = strings.ReplaceAll(jName, "/", "")
			jName = strings.ReplaceAll(jName, "\\", "")
			job.Name = &jName
		}
		job.Options.Timeout = loadTesterJob.Target.Timeout
		if job.Options.Timeout == 0 {
			job.Options.Timeout = 120000
		}

		if loadTesterJob.Target.LogResponse {
			job.Options.LogResult = true
		}

		if loadTesterJob.Target.Body != "" {
			job.Target.Body = loadTesterJob.Target.Body
		}
		if loadTesterJob.Target.ContentType != "" {
			job.Target.ContentType = loadTesterJob.Target.ContentType
		}
		if loadTesterJob.Target.BearerToken != "" {
			job.Target.JwtTokens = append(job.Target.JwtTokens, loadTesterJob.Target.BearerToken)
		}
		if loadTesterJob.Target.BearerTokens != nil && len(loadTesterJob.Target.BearerTokens) > 0 {
			for _, token := range loadTesterJob.Target.BearerTokens {
				job.Target.JwtTokens = append(job.Target.JwtTokens, token)
			}
		}
		if loadTesterJob.Target.Method != "" {
			job.Target.Method = job.Target.Method.Get(loadTesterJob.Target.Method)
		}
		if loadTesterJob.Target.URL != "" {
			job.Target.URL = loadTesterJob.Target.URL
		}
		if loadTesterJob.Target.Timeout > 0 {
			job.Options.Timeout = loadTesterJob.Target.Timeout
		}
		if loadTesterJob.Target.LogResponse {
			job.Target.logResponse = loadTesterJob.Target.LogResponse
		}

		if loadTesterJob.ConstantLoad != nil {
			job.Type = Constant

			switch strings.ToLower(loadTesterJob.ConstantLoad.Options.BlockType) {
			case "parallel":
				job.Options.BlockType = ParallelBlock
			case "sequential":
				job.Options.BlockType = SequentialBlock
			default:
				job.Options.BlockType = SequentialBlock
			}

			if loadTesterJob.ConstantLoad.Duration > 0 {
				job.Options.Duration = loadTesterJob.ConstantLoad.Duration
			}
			if loadTesterJob.ConstantLoad.Options.BlockInterval > 0 {
				job.Options.BlockInterval = NewInterval(loadTesterJob.ConstantLoad.Options.BlockInterval)
			}
			if loadTesterJob.ConstantLoad.Options.CallsPerBlock > 0 {
				job.Options.TasksPerBlock = NewInterval(loadTesterJob.ConstantLoad.Options.CallsPerBlock)
			}
			if loadTesterJob.ConstantLoad.Options.MinTaskInterval > 0 {
				job.Options.MinTaskInterval = NewInterval(loadTesterJob.ConstantLoad.Options.MinTaskInterval)
			}
			if loadTesterJob.ConstantLoad.Options.MaxTaskInterval > 0 {
				job.Options.MaxTaskInterval = NewInterval(loadTesterJob.ConstantLoad.Options.MaxTaskInterval)
			}
			if loadTest.Report.MaxTaskOutput > 0 {
				job.Options.MaxTaskOutput = loadTest.Report.MaxTaskOutput
			}
		} else if loadTesterJob.IncreasingLoad != nil {
			job.Type = Increasing

			switch strings.ToLower(loadTesterJob.IncreasingLoad.Options.BlockType) {
			case "parallel":
				job.Options.BlockType = ParallelBlock
			case "sequential":
				job.Options.BlockType = SequentialBlock
			default:
				job.Options.BlockType = SequentialBlock
			}

			if loadTesterJob.IncreasingLoad.Duration > 0 {
				job.Options.Duration = loadTesterJob.IncreasingLoad.Duration
			}
			if loadTesterJob.IncreasingLoad.Options.BlockInterval > 0 {
				job.Options.BlockInterval = NewInterval(loadTesterJob.IncreasingLoad.Options.BlockInterval)
			}
			if loadTesterJob.IncreasingLoad.Options.TotalCalls > 0 {
				job.Options.TasksPerBlock = NewInterval(loadTesterJob.IncreasingLoad.Options.TotalCalls)
			}
			if loadTesterJob.IncreasingLoad.Options.MinTaskInterval > 0 {
				job.Options.MinTaskInterval = NewInterval(loadTesterJob.IncreasingLoad.Options.MinTaskInterval)
			}
			if loadTesterJob.IncreasingLoad.Options.MaxTaskInterval > 0 {
				job.Options.MaxTaskInterval = NewInterval(loadTesterJob.IncreasingLoad.Options.MaxTaskInterval)
			}
			if loadTest.Report.MaxTaskOutput > 0 {
				job.Options.MaxTaskOutput = loadTest.Report.MaxTaskOutput
			}
		} else if loadTesterJob.FuzzyLoad != nil {
			job.Type = Fuzzy

			switch strings.ToLower(loadTesterJob.FuzzyLoad.Options.BlockType) {
			case "parallel":
				job.Options.BlockType = ParallelBlock
			case "sequential":
				job.Options.BlockType = SequentialBlock
			default:
				job.Options.BlockType = SequentialBlock
			}

			if loadTesterJob.FuzzyLoad.Duration > 0 {
				job.Options.Duration = loadTesterJob.FuzzyLoad.Duration
			}
			if loadTesterJob.FuzzyLoad.Options.MaxBlockInterval > 0 {
				job.Options.MaxBlockInterval = NewInterval(loadTesterJob.FuzzyLoad.Options.MaxBlockInterval)
			}
			if loadTesterJob.FuzzyLoad.Options.MinBlockInterval > 0 {
				job.Options.MinBlockInterval = NewInterval(loadTesterJob.FuzzyLoad.Options.MinBlockInterval)
			}
			if loadTesterJob.FuzzyLoad.Options.MaxTasksPerBlock > 0 {
				job.Options.MaxTasksPerBlock = NewInterval(loadTesterJob.FuzzyLoad.Options.MaxTasksPerBlock)
			}
			if loadTesterJob.FuzzyLoad.Options.MinTasksPerBlock > 0 {
				job.Options.MinTasksPerBlock = NewInterval(loadTesterJob.FuzzyLoad.Options.MinTasksPerBlock)
			}
			if loadTesterJob.FuzzyLoad.Options.MinTaskInterval > 0 {
				job.Options.MinTaskInterval = NewInterval(loadTesterJob.FuzzyLoad.Options.MinTaskInterval)
			}
			if loadTesterJob.FuzzyLoad.Options.MaxTaskInterval > 0 {
				job.Options.MaxTaskInterval = NewInterval(loadTesterJob.FuzzyLoad.Options.MaxTaskInterval)
			}
			if loadTest.Report.MaxTaskOutput > 0 {
				job.Options.MaxTaskOutput = loadTest.Report.MaxTaskOutput
			}
		}

		loadTesterJobs = append(loadTesterJobs, job)
	}
	logger.Success("Created successfully %v job instructions to execute", fmt.Sprintf("%v", len(loadTesterJobs)))

	var loadTestJobsWaitGroup sync.WaitGroup
	loadTestJobsWaitGroup.Add(len(loadTest.Jobs))

	for i, jobToExecute := range loadTesterJobs {
		switch strings.ToLower(loadTest.JobType) {
		case "parallel":
			logger.Command("Executing job %v in parallel", *jobToExecute.Name)
			go jobToExecute.Execute(&loadTestJobsWaitGroup)
		case "sequential":
			logger.Command("Executing job %v in sequence", *jobToExecute.Name)
			jobToExecute.Execute(&loadTestJobsWaitGroup)
		default:
			logger.Command("Executing job %v in sequence", *jobToExecute.Name)
			jobToExecute.Execute(&loadTestJobsWaitGroup)
		}

		if i < len(loadTest.Jobs)-1 {
			if loadTest.WaitBetweenJobs > 0 {
				logger.Info("Waiting for %v Millisecond for the next job...", fmt.Sprint(loadTest.WaitBetweenJobs))
				time.Sleep(time.Duration(loadTest.WaitBetweenJobs) * time.Millisecond)
			}
		}
	}

	loadTestJobsWaitGroup.Wait()

	return loadTesterJobs, nil
}
