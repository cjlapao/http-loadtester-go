package domain

import (
	"fmt"
	"sync"
	"time"
)

type VirtualUserJobOperation struct{}

func (v VirtualUserJobOperation) Execute(jobOperation *JobOperation) error {
	if jobOperation.Target.IsMultiTargeted() {
		logger.Success("Performing Virtual User Load Test on %v targets with %v virtual users", fmt.Sprintf("%v", jobOperation.Target.CountUrls()), fmt.Sprint(jobOperation.Options.NumberOfBlocks))
	} else {
		logger.Success("Performing Virtual User Load Test on %v with %v virtual users", jobOperation.Target.GetUrl(0), fmt.Sprint(jobOperation.Options.NumberOfBlocks))
	}

	startingJobTime := time.Now().UTC()
	jobOperation.Result = jobOperation.CreateLoadJobResult()

	GenerateVirtualUserBlocks(jobOperation, jobOperation.Options.NumberOfBlocks)
	amountOfBlocks := len(jobOperation.Blocks)
	var blockWaitingGroup sync.WaitGroup
	blockWaitingGroup.Add(amountOfBlocks)

	startingTime := time.Now().UTC()

	// Executing virtual users in parallel
	for i, block := range jobOperation.Blocks {
		blockNum := i + 1
		block.BlockPosition = blockNum
		block.TotalBlocks = amountOfBlocks
		logger.Info("Started processing virtual user %v [%v/%v] with %v timeout", block.ID, fmt.Sprint(blockNum), fmt.Sprint(amountOfBlocks), fmt.Sprint(time.Duration(jobOperation.Options.Timeout)*time.Millisecond))
		go block.ExecuteUntil(&blockWaitingGroup)
	}

	for {
		expectedEndingTime := startingJobTime.Add(time.Duration(jobOperation.Options.Duration) * time.Millisecond)
		shouldFinish := time.Now().UTC().After(expectedEndingTime)
		if shouldFinish {
			for _, block := range jobOperation.Blocks {
				close(block.StopChannel)
			}

			totalTestTime := time.Now().UTC().Sub(startingJobTime)
			if jobOperation.Target.IsMultiTargeted() {
				logger.Success("Finished Load Test on %v targets for %v", fmt.Sprintf("%v", jobOperation.Target.CountUrls()), totalTestTime.String())

			} else {
				logger.Success("Finished Load Test on %v for %v", jobOperation.Target.GetUrl(0), totalTestTime.String())
			}

			break
		}
	}

	blockWaitingGroup.Wait()

	endingTime := time.Now().UTC()
	var duration time.Duration = endingTime.Sub(startingTime)

	jobOperation.Result.TotalDurationInSeconds = duration.Seconds()
	// Parsing the results after we go all of them done
	for _, block := range jobOperation.Blocks {
		block.Result.Process()
		*jobOperation.Result.BlockResults = append(*jobOperation.Result.BlockResults, block.Result)
	}

	jobOperation.Result.ProcessResult()
	endingJobTime := time.Now().UTC()

	jobOperation.Result.TimeTaken = endingJobTime.Sub(startingJobTime)
	jobOperation.Result.StartingTime = startingJobTime
	jobOperation.Result.EndingTime = endingJobTime

	// if services.IsDatabaseEnabled() {
	// 	services.JobOperationRepo().Upsert(*jobOperation)
	// }

	return nil
}

func GenerateVirtualUserBlocks(operation *JobOperation, numberOfBlocks int) error {
	generateVirtualUserBlocks(operation, numberOfBlocks)
	return nil
}

// Generates virtual user job blocks
func generateVirtualUserBlocks(job *JobOperation, numberBlocks int) {
	for i := 0; i < numberBlocks; i++ {
		block := job.CreateBlock()
		block.BlockType = ParallelBlock
		block.TaskBlockType = SequentialBlock
		block.MinTaskInterval = job.Options.MinTaskInterval
		block.MaxTaskInterval = job.Options.MaxTaskInterval
		if job.Options.MaxBlockInterval.Value() > job.Options.MinBlockInterval.value {
			block.WaitFor = NewInterval(job.getRandomBlockInterval())
		} else {
			block.WaitFor = NewInterval(job.Options.MaxBlockInterval.Value())
		}

		block.CreateTask(0)
	}
}
