package jobs

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cjlapao/common-go/helper"
)

// JobOperationBlockTask Entity
type JobOperationBlockTask struct {
	ID              int
	BlockID         string
	JobID           string
	JobName         *string
	Timeout         int
	Target          *JobOperationTarget
	Type            JobOPerationType
	MinTaskInterval Interval
	MaxTaskInterval Interval
	Result          *JobOperationBlockTaskResult
	Verbose         bool
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
	if *j.JobName != "" {
		task.JobName = j.JobName
	} else {
		task.JobName = &j.JobID
	}

	task.Verbose = helper.GetFlagSwitch("verbose", false)
	task.MaxTaskInterval = j.MaxTaskInterval
	task.MinTaskInterval = j.MinTaskInterval

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

	if t.MaxTaskInterval.Value() > t.MinTaskInterval.Value() {
		waitFor := GetRandomBlockInterval(t.MaxTaskInterval, t.MinTaskInterval)
		if waitFor > 0 {
			time.Sleep(time.Duration(waitFor) * time.Millisecond)
		}
	}

	if t.Verbose {
		logger.Info("[%v] Started call %v to %v", *t.JobName, fmt.Sprint(t.ID), t.Target.URL)
	}

	// Implementing defined minutes timeout
	client := &http.Client{
		Timeout: time.Duration(t.Timeout) * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Millisecond)
	defer cancel()

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
			logger.Error("[%v] Error creating request %v on call %v: %v", *t.JobName, t.Target.Method.String(), fmt.Sprint(t.ID), err.Error())
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
	response, err = client.Do(request.WithContext(ctx))
	endingTime := time.Now().UTC()

	if err != nil {
		if t.Verbose {
			logger.Error("[%v] Error on call %v: %v", *t.JobName, fmt.Sprint(t.ID), err.Error())
		}

		var duration time.Duration = endingTime.Sub(startingTime)

		queryDuration := JobOperationBlockTaskDuration{
			Duration: duration,
			Seconds:  duration.Seconds(),
		}

		errorString := err.Error()
		if strings.Contains(errorString, "target machine actively refused it") {
			t.Result.QueryDuration = &queryDuration
			t.Result.StatusCode = 408
			t.Result.Status = "4xx Target machine actively refused it"
			t.Result.ErrorMessage = errorString
		} else {
			t.Result.QueryDuration = &queryDuration
			if response != nil {
				t.Result.StatusCode = response.StatusCode
				t.Result.Status = response.Status
			} else {
				t.Result.StatusCode = 999
				t.Result.Status = "xxx Empty Response"
			}
			t.Result.ErrorMessage = errorString
		}

		wg.Done()
		return
	}

	var duration time.Duration = endingTime.Sub(startingTime)
	if t.Verbose {
		logger.Info("[%v] Ended call %v to %v, took %v", *t.JobName, fmt.Sprint(t.ID), t.Target.URL, fmt.Sprint(duration.Seconds()))
	}

	queryDuration := JobOperationBlockTaskDuration{
		Duration: duration,
		Seconds:  duration.Seconds(),
	}

	t.Result.QueryDuration = &queryDuration

	responseDetails := ResponseDetails{}

	if response != nil {
		responseContent, err := ioutil.ReadAll(response.Body)
		if err != nil {
			if t.Verbose {
				logger.Error("[%v] Error reading content on call %v: %v", *t.JobName, fmt.Sprint(t.ID), err.Error())
			}
		}

		responseDetails.IP = getIP(response.Request)

		if response.TLS != nil {
			responseDetails.TLSCipher = tls.CipherSuiteName(response.TLS.CipherSuite)
			responseDetails.TLSServerName = response.TLS.ServerName

			switch response.TLS.Version {
			case 769:
				responseDetails.TLSVersion = "TLSv1.0"
			case 770:
				responseDetails.TLSVersion = "TLSv1.1"
			case 771:
				responseDetails.TLSVersion = "TLSv1.2"
			case 772:
				responseDetails.TLSVersion = "TLSv1.3"
			}
		}

		t.Result.ResponseDetails = &responseDetails
		t.Result.StatusCode = response.StatusCode
		t.Result.Status = response.Status
		t.Result.Content = string(responseContent)

	} else {
		logger.Error("There was no response from the api")
		t.Result.ResponseDetails = &responseDetails
		t.Result.Status = "No Response back from server"
		t.Result.StatusCode = 500
	}

	wg.Done()
}

func getIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// JobOperationBlockTaskResult Entity
type JobOperationBlockTaskResult struct {
	TaskID          int
	BlockID         string
	JobID           string
	Target          *JobOperationTarget
	QueryDuration   *JobOperationBlockTaskDuration
	Status          string
	StatusCode      int
	Content         string
	ErrorMessage    string
	ResponseDetails *ResponseDetails
}

// JobOperationBlockTaskDuration Entity
type JobOperationBlockTaskDuration struct {
	Duration     time.Duration
	Milliseconds int64
	Seconds      float64
}
