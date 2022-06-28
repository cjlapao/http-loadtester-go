package usecases

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/http-loadtester-go/domain"
	"github.com/cjlapao/http-loadtester-go/entities"
	"github.com/cjlapao/http-loadtester-go/infrastructure"

	"gopkg.in/yaml.v2"
)

// ExecuteFromFile Execute LoadTest from file
func ExecuteFromFile(filepath string) error {
	if !helper.FileExists(filepath) {
		err := errors.New("file not found")
		sp.Logger().Error(err.Error())
		return err
	}

	content, err := helper.ReadFromFile(filepath)
	if err != nil {
		sp.Logger().Error(err.Error())
		return err
	}

	var loadTest entities.LoadTest
	err = yaml.Unmarshal(content, &loadTest)
	if err != nil {
		sp.Logger().Error(err.Error())
		return err
	}

	jobs, err := ExecuteLoadTest(loadTest)
	if err != nil {
		return err
	}

	if loadTest.Report.OutputToFile {
		for _, job := range jobs {
			job.ExportReportToFile(loadTest.Report.OutputToFilePath)
			sp.Logger().Success("Finished creating reports for job %v", *job.Name)
		}
	} else {
		for _, job := range jobs {
			fmt.Println(job.MarkDown())
			sp.Logger().Success("Finished creating reports for job %v", *job.Name)
		}
	}
	if loadTest.Report.OutputResults {
		for _, job := range jobs {

			job.ExportOutputToFile(loadTest.Report.OutputToFilePath)
			sp.Logger().Success("Finished creating reports for job %v", *job.Name)
		}
	}

	return nil
}

func ExecuteLoadTest(loadTest entities.LoadTest) ([]*domain.JobOperation, error) {
	if loadTest.DisplayName != "" {
		sp.Logger().Success("Starting %v load test.", loadTest.DisplayName)
	}
	loadTesterJobs := make([]*domain.JobOperation, 0)

	for _, loadTesterJob := range loadTest.Jobs {
		if loadTesterJob.Disabled {
			continue
		}

		if loadTesterJob.Target.URL == "" && (loadTesterJob.Target.URLs == nil || len(loadTesterJob.Target.URLs) == 0) {
			err := errors.New("url was not defined")
			sp.Logger().Error(err.Error())
			return loadTesterJobs, err
		}

		job := domain.CreateJobOperation(infrastructure.GetServiceProvider())

		switch strings.ToLower(loadTesterJob.Type) {
		case "parallel":
			job.OperationType = domain.ParallelBlock
		case "sequential":
			job.OperationType = domain.SequentialBlock
		default:
			job.OperationType = domain.SequentialBlock
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

		// Bearer Tokens
		if loadTesterJob.Target.BearerToken != "" {
			job.Target.JwtTokens = append(job.Target.JwtTokens, loadTesterJob.Target.BearerToken)
		}
		if loadTesterJob.Target.BearerTokens != nil && len(loadTesterJob.Target.BearerTokens) > 0 {
			job.Target.JwtTokens = append(job.Target.JwtTokens, loadTesterJob.Target.BearerTokens...)
		}

		// Method
		if loadTesterJob.Target.Method != "" {
			job.Target.Method = job.Target.Method.Get(loadTesterJob.Target.Method)
		}

		// Single targeted url
		if loadTesterJob.Target.URL != "" {
			job.Target.URLs = append(job.Target.URLs, loadTesterJob.Target.URL)
		}

		// Single targeted url
		if loadTesterJob.Target.URLs != nil && len(loadTesterJob.Target.URLs) > 0 {
			job.Target.URLs = append(job.Target.URLs, loadTesterJob.Target.URLs...)
		}

		// UserAgent
		if loadTesterJob.Target.UserAgent != "" {
			job.Target.UserAgent = loadTesterJob.Target.UserAgent
		}

		// Headers
		if loadTesterJob.Target.Headers != nil && len(loadTesterJob.Target.Headers) > 0 {
			for key, value := range loadTesterJob.Target.Headers {
				job.Target.Headers[key] = value
			}
		}

		if loadTesterJob.Target.Timeout > 0 {
			job.Options.Timeout = loadTesterJob.Target.Timeout
		}
		if loadTesterJob.Target.LogResponse {
			job.Target.LogResponse = loadTesterJob.Target.LogResponse
		}

		if loadTesterJob.ConstantLoad != nil {
			job.Type = domain.Constant

			switch strings.ToLower(loadTesterJob.ConstantLoad.Options.BlockType) {
			case "parallel":
				job.Options.BlockType = domain.ParallelBlock
			case "sequential":
				job.Options.BlockType = domain.SequentialBlock
			default:
				job.Options.BlockType = domain.SequentialBlock
			}

			if loadTesterJob.ConstantLoad.Duration > 0 {
				job.Options.Duration = loadTesterJob.ConstantLoad.Duration
			}
			if loadTesterJob.ConstantLoad.Options.NumberOfBlocks > 0 {
				job.Options.NumberOfBlocks = loadTesterJob.ConstantLoad.Options.NumberOfBlocks
			}
			if loadTesterJob.ConstantLoad.Options.BlockInterval > 0 {
				job.Options.BlockInterval = domain.NewInterval(loadTesterJob.ConstantLoad.Options.BlockInterval)
			}
			if loadTesterJob.ConstantLoad.Options.CallsPerBlock > 0 {
				job.Options.TasksPerBlock = domain.NewInterval(loadTesterJob.ConstantLoad.Options.CallsPerBlock)
			}
			if loadTesterJob.ConstantLoad.Options.MinTaskInterval > 0 {
				job.Options.MinTaskInterval = domain.NewInterval(loadTesterJob.ConstantLoad.Options.MinTaskInterval)
			}
			if loadTesterJob.ConstantLoad.Options.MaxTaskInterval > 0 {
				job.Options.MaxTaskInterval = domain.NewInterval(loadTesterJob.ConstantLoad.Options.MaxTaskInterval)
			}
			if loadTest.Report.MaxTaskOutput > 0 {
				job.Options.MaxTaskOutput = loadTest.Report.MaxTaskOutput
			}
		} else if loadTesterJob.IncreasingLoad != nil {
			job.Type = domain.Increasing

			switch strings.ToLower(loadTesterJob.IncreasingLoad.Options.BlockType) {
			case "parallel":
				job.Options.BlockType = domain.ParallelBlock
			case "sequential":
				job.Options.BlockType = domain.SequentialBlock
			default:
				job.Options.BlockType = domain.SequentialBlock
			}

			if loadTesterJob.IncreasingLoad.Duration > 0 {
				job.Options.Duration = loadTesterJob.IncreasingLoad.Duration
			}
			if loadTesterJob.IncreasingLoad.Options.NumberOfBlocks > 0 {
				job.Options.NumberOfBlocks = loadTesterJob.IncreasingLoad.Options.NumberOfBlocks
			}
			if loadTesterJob.IncreasingLoad.Options.BlockInterval > 0 {
				job.Options.BlockInterval = domain.NewInterval(loadTesterJob.IncreasingLoad.Options.BlockInterval)
			}
			if loadTesterJob.IncreasingLoad.Options.TotalCalls > 0 {
				job.Options.TasksPerBlock = domain.NewInterval(loadTesterJob.IncreasingLoad.Options.TotalCalls)
			}
			if loadTesterJob.IncreasingLoad.Options.MinTaskInterval > 0 {
				job.Options.MinTaskInterval = domain.NewInterval(loadTesterJob.IncreasingLoad.Options.MinTaskInterval)
			}
			if loadTesterJob.IncreasingLoad.Options.MaxTaskInterval > 0 {
				job.Options.MaxTaskInterval = domain.NewInterval(loadTesterJob.IncreasingLoad.Options.MaxTaskInterval)
			}
			if loadTest.Report.MaxTaskOutput > 0 {
				job.Options.MaxTaskOutput = loadTest.Report.MaxTaskOutput
			}
		} else if loadTesterJob.FuzzyLoad != nil {
			job.Type = domain.Fuzzy

			switch strings.ToLower(loadTesterJob.FuzzyLoad.Options.BlockType) {
			case "parallel":
				job.Options.BlockType = domain.ParallelBlock
			case "sequential":
				job.Options.BlockType = domain.SequentialBlock
			default:
				job.Options.BlockType = domain.SequentialBlock
			}

			if loadTesterJob.FuzzyLoad.Duration > 0 {
				job.Options.Duration = loadTesterJob.FuzzyLoad.Duration
			}
			if loadTesterJob.FuzzyLoad.Options.NumberOfBlocks > 0 {
				job.Options.NumberOfBlocks = loadTesterJob.FuzzyLoad.Options.NumberOfBlocks
			}
			if loadTesterJob.FuzzyLoad.Options.MaxBlockInterval > 0 {
				job.Options.MaxBlockInterval = domain.NewInterval(loadTesterJob.FuzzyLoad.Options.MaxBlockInterval)
			}
			if loadTesterJob.FuzzyLoad.Options.MinBlockInterval > 0 {
				job.Options.MinBlockInterval = domain.NewInterval(loadTesterJob.FuzzyLoad.Options.MinBlockInterval)
			}
			if loadTesterJob.FuzzyLoad.Options.MaxTasksPerBlock > 0 {
				job.Options.MaxTasksPerBlock = domain.NewInterval(loadTesterJob.FuzzyLoad.Options.MaxTasksPerBlock)
			}
			if loadTesterJob.FuzzyLoad.Options.MinTasksPerBlock > 0 {
				job.Options.MinTasksPerBlock = domain.NewInterval(loadTesterJob.FuzzyLoad.Options.MinTasksPerBlock)
			}
			if loadTesterJob.FuzzyLoad.Options.MinTaskInterval > 0 {
				job.Options.MinTaskInterval = domain.NewInterval(loadTesterJob.FuzzyLoad.Options.MinTaskInterval)
			}
			if loadTesterJob.FuzzyLoad.Options.MaxTaskInterval > 0 {
				job.Options.MaxTaskInterval = domain.NewInterval(loadTesterJob.FuzzyLoad.Options.MaxTaskInterval)
			}
			if loadTest.Report.MaxTaskOutput > 0 {
				job.Options.MaxTaskOutput = loadTest.Report.MaxTaskOutput
			}
		}

		loadTesterJobs = append(loadTesterJobs, job)
	}

	sp.Logger().Success("Created successfully %v job instructions to execute", fmt.Sprintf("%v", len(loadTesterJobs)))

	var loadTestJobsWaitGroup sync.WaitGroup
	loadTestJobsWaitGroup.Add(len(loadTesterJobs))

	for i, jobToExecute := range loadTesterJobs {
		switch strings.ToLower(loadTest.JobType) {
		case "parallel":
			sp.Logger().Command("Executing job %v in parallel", *jobToExecute.Name)
			go jobToExecute.Execute(&loadTestJobsWaitGroup)
		case "sequential":
			sp.Logger().Command("Executing job %v in sequence", *jobToExecute.Name)
			jobToExecute.Execute(&loadTestJobsWaitGroup)
		default:
			sp.Logger().Command("Executing job %v in sequence", *jobToExecute.Name)
			jobToExecute.Execute(&loadTestJobsWaitGroup)
		}

		if i < len(loadTesterJobs)-1 {
			if loadTest.WaitBetweenJobs > 0 {
				sp.Logger().Info("Waiting for %v Millisecond for the next job...", fmt.Sprint(loadTest.WaitBetweenJobs))
				time.Sleep(time.Duration(loadTest.WaitBetweenJobs) * time.Millisecond)
			}
		}
	}

	loadTestJobsWaitGroup.Wait()

	return loadTesterJobs, nil
}
