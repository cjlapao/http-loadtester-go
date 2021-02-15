package jobs

import (
	"errors"
	"fmt"
	"strings"
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
	URL              string `json:"url" yaml:"url"`
	Method           string `json:"method" yaml:"method"`
	Body             string `json:"body" yaml:"body"`
	BearerToken      string `json:"token" yaml:"token"`
	ContentType      string `json:"contentType" yaml:"contentType"`
	TimeoutInSeconds int    `json:"timeout" yaml:"timeout"`
}

// LoadTestConstantJob entity
type LoadTestConstantJob struct {
	Duration int                        `json:"duration" yaml:"duration"`
	Options  LoadTestConstantJobOptions `json:"specs" yaml:"specs"`
}

// LoadTestConstantJobOptions Entity
type LoadTestConstantJobOptions struct {
	BlockInterval int `json:"blockInterval" yaml:"blockInterval"`
	CallsPerBlock int `json:"callsPerBlock" yaml:"callsPerBlock"`
}

// LoadTestIncreasingJob Entity
type LoadTestIncreasingJob struct {
	Duration int                          `json:"duration" yaml:"duration"`
	Options  LoadTestIncreasingJobOptions `json:"specs" yaml:"specs"`
}

// LoadTestIncreasingJobOptions Entity
type LoadTestIncreasingJobOptions struct {
	BlockInterval int `json:"blockInterval" yaml:"blockInterval"`
	TotalCalls    int `json:"totalCalls" yaml:"totalCalls"`
}

// LoadTestFuzzyJob Entity
type LoadTestFuzzyJob struct {
	Duration int                     `json:"duration" yaml:"duration"`
	Options  LoadTestFuzzyJobOptions `json:"specs" yaml:"specs"`
}

// LoadTestFuzzyJobOptions Entity
type LoadTestFuzzyJobOptions struct {
	MaxBlockInterval int `json:"maxBlockInterval" yaml:"maxBlockInterval"`
	MinBlockInterval int `json:"minBlockInterval" yaml:"minBlockInterval"`
	MaxTasksPerBlock int `json:"maxTaskPerBlock" yaml:"maxTaskPerBlock"`
	MinTasksPerBlock int `json:"minTaskPerBlock" yaml:"minTaskPerBlock"`
}

// LoadTestJobOutput Entity
type LoadTestJobOutput struct {
	MaxTaskOutput int  `json:"maxTaskOutput" yaml:"maxTaskOutput"`
	OutputResults bool `json:"outputResults" yaml:"outputResults"`
	OutputToFile  bool `json:"outputToFile" yaml:"outputToFile"`
}

// ExecuteFromFile Execute LoadTest from file
func ExecuteFromFile(filepath string) error {
	if !helper.FileExists(filepath) {
		err := errors.New("File was not found")
		logger.LogError(err)
		return err
	}

	content, err := helper.ReadFromFile(filepath)
	if err != nil {
		logger.LogError(err)
		return err
	}

	var loadTest LoadTest
	err = yaml.Unmarshal(content, &loadTest)
	if err != nil {
		logger.LogError(err)
		return err
	}
	if loadTest.DisplayName != "" {
		logger.Success("Testing File %v ", loadTest.DisplayName)
	}

	for i, loadTesterJob := range loadTest.Jobs {
		if loadTesterJob.Target.URL == "" {
			err := errors.New("Url was not defined")
			logger.LogError(err)
			return err
		}
		job := CreateJobOperation()

		switch strings.ToLower(loadTesterJob.Type) {
		case "parallel":
			job.BlockType = ParallelBlock
		case "sequential":
			job.BlockType = SequentialBlock
		default:
			job.BlockType = ParallelBlock
		}
		if loadTesterJob.Name != "" {
			job.Name = strings.ReplaceAll(loadTesterJob.Name, " ", "_")
			job.Name = strings.ReplaceAll(job.Name, ":", "")
			job.Name = strings.ReplaceAll(job.Name, "/", "")
			job.Name = strings.ReplaceAll(job.Name, "\\", "")
		}
		if loadTesterJob.Target.Body != "" {
			job.Target.Body = loadTesterJob.Target.Body
		}
		if loadTesterJob.Target.ContentType != "" {
			job.Target.ContentType = loadTesterJob.Target.ContentType
		}
		if loadTesterJob.Target.BearerToken != "" {
			job.Target.JwtToken = loadTesterJob.Target.BearerToken
		}
		if loadTesterJob.Target.Method != "" {
			job.Target.Method = job.Target.Method.Get(loadTesterJob.Target.Method)
		}
		if loadTesterJob.Target.URL != "" {
			job.Target.URL = loadTesterJob.Target.URL
		}
		if loadTesterJob.Target.TimeoutInSeconds > 0 {
			job.Options.Timeout = loadTesterJob.Target.TimeoutInSeconds
		}
		if loadTesterJob.ConstantLoad != nil {
			job.Type = Constant
			if loadTesterJob.ConstantLoad.Duration > 0 {
				job.Options.Duration = loadTesterJob.ConstantLoad.Duration
			}
			if loadTesterJob.ConstantLoad.Options.BlockInterval > 0 {
				job.Options.BlockInterval = NewInterval(loadTesterJob.ConstantLoad.Options.BlockInterval)
			}
			if loadTesterJob.ConstantLoad.Options.CallsPerBlock > 0 {
				job.Options.TasksPerBlock = NewInterval(loadTesterJob.ConstantLoad.Options.CallsPerBlock)
			}
			if loadTest.Report.MaxTaskOutput > 0 {
				job.Options.MaxTaskOutput = loadTest.Report.MaxTaskOutput
			}
		} else if loadTesterJob.IncreasingLoad != nil {
			job.Type = Increasing
			if loadTesterJob.IncreasingLoad.Duration > 0 {
				job.Options.Duration = loadTesterJob.IncreasingLoad.Duration
			}
			if loadTesterJob.IncreasingLoad.Options.BlockInterval > 0 {
				job.Options.BlockInterval = NewInterval(loadTesterJob.IncreasingLoad.Options.BlockInterval)
			}
			if loadTesterJob.IncreasingLoad.Options.TotalCalls > 0 {
				job.Options.TasksPerBlock = NewInterval(loadTesterJob.IncreasingLoad.Options.TotalCalls)
			}
			if loadTest.Report.MaxTaskOutput > 0 {
				job.Options.MaxTaskOutput = loadTest.Report.MaxTaskOutput
			}
		} else if loadTesterJob.FuzzyLoad != nil {
			job.Type = Fuzzy
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
			if loadTest.Report.MaxTaskOutput > 0 {
				job.Options.MaxTaskOutput = loadTest.Report.MaxTaskOutput
			}
		}
		logger.Success("Starting job %v execution.", job.Name)
		job.Execute()
		logger.Success("Finished executing job %v, generating reports...", job.Name)
		if loadTest.Report.OutputToFile {
			job.ExportReportToFile()
		} else {
			fmt.Println(job.MarkDown())
		}
		if loadTest.Report.OutputResults {
			job.ExportOutputToFile()
		}
		logger.Success("Finished creating reports for job %v", job.Name)
		if i < len(loadTest.Jobs)-1 {
			if loadTest.WaitBetweenJobs > 0 {
				logger.Info("Waiting for %v seconds for the next job...", fmt.Sprint(loadTest.WaitBetweenJobs))
				time.Sleep(time.Duration(loadTest.WaitBetweenJobs) * time.Second)
			}
		}
	}
	return nil
}
