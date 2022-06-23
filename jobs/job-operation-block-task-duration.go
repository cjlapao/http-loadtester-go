package jobs

import "time"

// JobOperationBlockTaskDuration Entity
type JobOperationBlockTaskDuration struct {
	Duration     time.Duration
	Milliseconds int64
	Seconds      float64
}
