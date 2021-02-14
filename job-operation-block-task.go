package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/cjlapao/common-go/helper"
)

// JobOperationBlockTask Entity
type JobOperationBlockTask struct {
	ID      int
	BlockID string
	JobID   string
	JobName string
	Timeout int
	Target  *JobOperationTarget
	Type    JobOPerationType
	Result  *JobOperationBlockTaskResult
	Verbose bool
}

// CreateTask Creates a Task Inside a Job Block
func (j *JobOperationBlock) CreateTask(id int) *JobOperationBlockTask {
	task := JobOperationBlockTask{
		ID:      id,
		BlockID: j.ID,
		JobID:   j.JobID,
		Target:  j.Target,
		Type:    j.Type,
		Timeout: j.Timeout,
	}
	if j.JobName != "" {
		task.JobName = j.JobName
	} else {
		task.JobName = j.JobID
	}

	task.Verbose = helper.GetFlagSwitch("verbose", false)

	if j.Tasks == nil {
		tasks := make([]*JobOperationBlockTask, 0)
		j.Tasks = &tasks
	}
	task.Result = task.CreateResult()
	*j.Tasks = append(*j.Tasks, &task)
	return &task
}

// CreateResult Creates a JobOperationBlockTask result entity
func (t *JobOperationBlockTask) CreateResult() *JobOperationBlockTaskResult {
	result := JobOperationBlockTaskResult{
		TaskID:     t.ID,
		BlockID:    t.BlockID,
		JobID:      t.JobID,
		Target:     t.Target,
		StatusCode: 100,
	}

	return &result
}

// Execute Executes a Sync Block Task
func (t *JobOperationBlockTask) Execute(wg *sync.WaitGroup) {
	if t.Target.URL == "" {
		wg.Done()
		return
	}

	if t.Verbose {
		logger.Info("[%v] Started call %v to %v", t.JobName, fmt.Sprint(t.ID), t.Target.URL)
	}

	// Implementing defined minutes timeout
	client := &http.Client{
		Timeout: time.Duration(t.Timeout) * time.Second,
	}

	var response *http.Response
	var request *http.Request
	var err error

	startingTime := time.Now().UTC()

	if t.Target.Body != "" {
		request, err = http.NewRequest(t.Target.Method.String(), t.Target.URL, bytes.NewReader([]byte(t.Target.Body)))
	} else {
		request, err = http.NewRequest(t.Target.Method.String(), t.Target.URL, nil)
	}

	if err != nil {
		if t.Verbose {
			logger.Error("[%v] Error creating request %v on call %v: %v", t.JobName, t.Target.Method.String(), fmt.Sprint(t.ID), err.Error())
		}
		wg.Done()
		return
	}

	if request != nil {
		if t.Target.JwtToken != "" {
			request.Header.Set("Authorization", "Bearer "+t.Target.JwtToken)
		}
		if t.Target.ContentType != "" {
			request.Header.Set("Content-Type", t.Target.ContentType)
		}
	}
	response, err = client.Do(request)
	endingTime := time.Now().UTC()
	if err != nil {
		if t.Verbose {
			logger.Error("[%v] Error on call %v: %v", t.JobName, fmt.Sprint(t.ID), err.Error())
		}

		var duration time.Duration = endingTime.Sub(startingTime)

		queryDuration := JobOperationBlockTaskDuration{
			Duration: duration,
			Seconds:  duration.Seconds(),
		}

		t.Result.QueryDuration = &queryDuration
		t.Result.StatusCode = 500
		t.Result.Status = "500 Load Balancer exception"

		wg.Done()
		return
	}

	var duration time.Duration = endingTime.Sub(startingTime)
	if t.Verbose {
		logger.Info("[%v] Ended cal %v to %v, took %v", t.JobName, fmt.Sprint(t.ID), t.Target.URL, fmt.Sprint(duration.Seconds()))
	}

	queryDuration := JobOperationBlockTaskDuration{
		Duration: duration,
		Seconds:  duration.Seconds(),
	}

	t.Result.QueryDuration = &queryDuration

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		if t.Verbose {
			logger.Error("[%v] Error reading content on call %v: %v", t.JobName, fmt.Sprint(t.ID), err.Error())
		}
	}

	t.Result.StatusCode = response.StatusCode
	t.Result.Status = response.Status
	t.Result.Content = string(responseContent)

	wg.Done()
}

// JobOperationBlockTaskResult Entity
type JobOperationBlockTaskResult struct {
	TaskID        int
	BlockID       string
	JobID         string
	Target        *JobOperationTarget
	QueryDuration *JobOperationBlockTaskDuration
	Status        string
	StatusCode    int
	Content       string
}

// JobOperationBlockTaskDuration Entity
type JobOperationBlockTaskDuration struct {
	Duration     time.Duration
	Milliseconds int64
	Seconds      float64
}
