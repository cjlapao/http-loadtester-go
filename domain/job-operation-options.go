package domain

// JobOperationOptions Entity
type JobOperationOptions struct {
	Duration         int
	NumberOfBlocks   int
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
