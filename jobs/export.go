package jobs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cjlapao/common-go/helper"

	"github.com/cjlapao/markdown-go"
	"gopkg.in/yaml.v2"
)

const (
	DDMMYYYYhhmmss = "2006_01_02_T15_04_05"
)

// ExportReportToFile Exports a Job Report to a File
func (j *JobOperation) ExportReportToFile(path string) {
	content := j.MarkDown()
	if path == "" {
		path = "./"
	} else {
		if !strings.HasSuffix(path, "/") {
			path = fmt.Sprintf("%v/", path)
		}
	}

	name := *j.Name + "_" + time.Now().Format(DDMMYYYYhhmmss)
	if name != j.ID {
		name += "_" + j.ID
	}

	fullPath := fmt.Sprintf("%vjob_%v_output.md", path, name)
	helper.WriteToFile(content, fullPath)
}

// ExportOutputToFile Exports a job result to a file
func (j *JobOperation) ExportOutputToFile(path string) {
	stringContent, err := yaml.Marshal(j.Result)
	if err != nil {
		return
	}
	if path == "" {
		path = "./"
	} else {
		if !strings.HasSuffix(path, "/") {
			path = fmt.Sprintf("%v/", path)
		}
	}

	name := *j.Name + "_" + time.Now().Format(DDMMYYYYhhmmss)
	if name != j.ID {
		name += "_" + j.ID
	}

	fullPath := fmt.Sprintf("%vjob_%v_output.yml", path, name)
	helper.WriteToFile(string(stringContent), fullPath)
}

// MarkDown Generates a Job Markdown report
func (j *JobOperation) MarkDown() string {
	md := markdown.CreateDocument()
	header := md.CreateHeader()
	if j.Target.IsMultiTargeted() {
		header.H1(fmt.Sprintf("HTTP Load Tester Report for %v targets", j.Target.CountUrls()))
	} else {
		header.H1(fmt.Sprintf("HTTP Load Tester Report for %v", j.Target.GetUrl(0)))
	}

	htb := md.CreateTextBlock()
	htb.AddLine(fmt.Sprintf("Job Type: %v", j.Type))
	htb.AddLine(fmt.Sprintf("Block Type: %v", j.OperationType))
	htb.AddLine(fmt.Sprintf("Task Type: %v", j.Options.BlockType))
	htb.AddLine(fmt.Sprintf("Timeout: %v", fmt.Sprint(time.Duration(j.Options.Timeout)*time.Millisecond)))
	htb.AddLine(fmt.Sprintf("Method: %v", j.Target.Method))
	if j.Target.Body != "" {
		htb.AddLine("Contains Body: Yes")
	} else {
		htb.AddLine("Contains Body: No")
	}
	htb.AddLine(fmt.Sprintf("Started: %v", j.Result.StartingTime.Format(time.RFC3339)))
	htb.AddLine(fmt.Sprintf("Finished: %v", j.Result.EndingTime.Format(time.RFC3339)))

	htb.AddLine(fmt.Sprintf("Duration: %.2f seconds", j.Result.TotalDurationInSeconds))
	htb.AddLine(fmt.Sprintf("Total Calls: %v, Succeeded: %v, Failed: %v", j.Result.TotalCalls, j.Result.TotalSucceededCalls, j.Result.TotalFailedCalls))
	if j.Result.TotalFailedCalls > 0 {
		percent := (float64(j.Result.TotalFailedCalls) * 100) / float64(j.Result.TotalCalls)
		htb.AddLine(fmt.Sprintf("Percent Failed: %.1f%%", percent))
	}
	htb.AddLine(fmt.Sprintf("Time Taken: %.2f seconds", j.Result.TimeTaken.Seconds()))
	htb.AddLine(fmt.Sprintf("Average Block Duration: %.4f seconds", j.Result.AverageBlockDuration))
	htb.AddLine(fmt.Sprintf("Average Call Duration: %.4f seconds", j.Result.AverageCallDuration))
	htb.AddLine(fmt.Sprintf("Authentication: %v", strconv.FormatBool(j.Authenticated())))
	htb.AddLine(fmt.Sprintf("Time Taken: %.2f seconds", j.Result.TimeTaken.Seconds()))
	htb.NewLine().NewLine()
	statusDetailsHeader := md.CreateHeader()
	statusDetailsHeader.H2("Status Results Details")
	j.generateTaskTable(md)
	if j.Target.IsMultiTargeted() {
		multiTargetDetails := md.CreateHeader()
		multiTargetDetails.H2("Targets")
		j.generateMultiTargetTable(md)
	}
	blockDetailsHeader := md.CreateHeader()
	blockDetailsHeader.H2("Block Results Details")
	j.generateBlockTable(md)
	blockTaskDetailsHeader := md.CreateHeader()
	blockTaskDetailsHeader.H2("Block Task Results Details")
	j.generateBlockTaskTable(md)

	return md.Sprint()
}

func (j *JobOperation) generateTaskTable(document *markdown.Document) *markdown.Table {
	table := document.CreateTable()
	table.AddAlignedHeaderColumn("Code", markdown.AlignRight)
	table.AddAlignedHeaderColumn("Total", markdown.AlignRight)

	for _, status := range *j.Result.TaskResponseStatus {
		table.AddRow(
			fmt.Sprintf("%v", status.Code),
			fmt.Sprintf("%v", status.Count),
		)
	}

	return table
}

func (j *JobOperation) generateBlockTable(document *markdown.Document) *markdown.Table {
	table := document.CreateTable()
	table.AddHeaderColumn("Block")
	table.AddAlignedHeaderColumn("Total Queries", markdown.AlignRight)
	table.AddAlignedHeaderColumn("Failed Queries", markdown.AlignRight)
	table.AddAlignedHeaderColumn("Succeeded Queries", markdown.AlignRight)
	table.AddAlignedHeaderColumn("Duration", markdown.AlignRight)
	table.AddAlignedHeaderColumn("Average Task Duration", markdown.AlignRight)
	for _, block := range *j.Result.BlockResults {
		table.AddRow(
			fmt.Sprintf("Block %v", block.BlockID),
			fmt.Sprint(block.Total),
			fmt.Sprint(block.Failed),
			fmt.Sprint(block.Succeeded),
			fmt.Sprintf("%.3f seconds", block.TotalDurationInSeconds),
			fmt.Sprintf("%.3f seconds", block.AverageTaskDuration),
		)
	}

	return table
}

func (j *JobOperation) generateMultiTargetTable(document *markdown.Document) *markdown.Table {
	table := document.CreateTable()
	table.AddHeaderColumn("Uri")
	table.AddAlignedHeaderColumn("Number of calls", markdown.AlignRight)
	for key, value := range j.Result.TargetCalls {
		table.AddRow(
			fmt.Sprint(key),
			fmt.Sprint(value),
		)
	}

	return table
}

func (j *JobOperation) generateBlockTaskTable(document *markdown.Document) *markdown.Table {
	table := document.CreateTable()
	table.AddHeaderColumn("Block")
	table.AddHeaderColumn("Tasks")
	table.AddAlignedHeaderColumn("Duration", markdown.AlignRight)
	table.AddAlignedHeaderColumn("Status Code", markdown.AlignCenter)
	if j.Options.LogResult {
		table.AddAlignedHeaderColumn("Content", markdown.AlignLeft)
	}
	for _, block := range *j.Result.BlockResults {
		table.AddRow(fmt.Sprintf("Start of Block %v", block.BlockID), "", "", "")
		bTask := *block.TaskResults
		if len(bTask) > j.Result.MaxTaskOutput {
			halfItems := j.Result.MaxTaskOutput / 2
			itemBlocks := j.Result.MaxTaskOutput / 3
			secondBlockStart := halfItems - (itemBlocks / 2)
			secondBlockEnd := halfItems + (itemBlocks / 2) + 1
			firstHalf := bTask[0:itemBlocks]
			for _, task := range firstHalf {
				if j.Options.LogResult {
					table.AddRow(
						fmt.Sprintf("Block %v", task.BlockID),
						fmt.Sprintf("Task %v", task.TaskID),
						fmt.Sprintf("%.2f seconds", task.QueryDuration.Seconds),
						fmt.Sprintf("%v", task.Status),
						fmt.Sprintf("%v", task.Content),
					)
				} else {
					table.AddRow(
						fmt.Sprintf("Block %v", task.BlockID),
						fmt.Sprintf("Task %v", task.TaskID),
						fmt.Sprintf("%.2f seconds", task.QueryDuration.Seconds),
						fmt.Sprintf("%v", task.Status),
					)

				}
			}
			if secondBlockStart < secondBlockEnd {
				table.AddRow("...", "...", "...", "...")
				secondHalf := bTask[secondBlockStart:secondBlockEnd]
				for _, task := range secondHalf {
					if j.Options.LogResult {
						table.AddRow(
							fmt.Sprintf("Block %v", task.BlockID),
							fmt.Sprintf("Task %v", task.TaskID),
							fmt.Sprintf("%.2f seconds", task.QueryDuration.Seconds),
							fmt.Sprintf("%v", task.Status),
							fmt.Sprintf("%v", task.Content),
						)
					} else {
						table.AddRow(
							fmt.Sprintf("Block %v", task.BlockID),
							fmt.Sprintf("Task %v", task.TaskID),
							fmt.Sprintf("%.2f seconds", task.QueryDuration.Seconds),
							fmt.Sprintf("%v", task.Status),
						)

					}
				}
			}
			table.AddRow("...", "...", "...", "...")
			thirdHalf := bTask[len(bTask)-itemBlocks:]
			for _, task := range thirdHalf {
				if j.Options.LogResult {
					table.AddRow(
						fmt.Sprintf("Block %v", task.BlockID),
						fmt.Sprintf("Task %v", task.TaskID),
						fmt.Sprintf("%.2f seconds", task.QueryDuration.Seconds),
						fmt.Sprintf("%v", task.Status),
						fmt.Sprintf("%v", task.Content),
					)
				} else {
					table.AddRow(
						fmt.Sprintf("Block %v", task.BlockID),
						fmt.Sprintf("Task %v", task.TaskID),
						fmt.Sprintf("%.2f seconds", task.QueryDuration.Seconds),
						fmt.Sprintf("%v", task.Status),
					)

				}
			}
		} else {
			for _, task := range bTask {
				if j.Options.LogResult {
					table.AddRow(
						fmt.Sprintf("Block %v", task.BlockID),
						fmt.Sprintf("Task %v", task.TaskID),
						fmt.Sprintf("%.2f seconds", task.QueryDuration.Seconds),
						fmt.Sprintf("%v", task.Status),
						fmt.Sprintf("%v", task.Content),
					)
				} else {
					table.AddRow(
						fmt.Sprintf("Block %v", task.BlockID),
						fmt.Sprintf("Task %v", task.TaskID),
						fmt.Sprintf("%.2f seconds", task.QueryDuration.Seconds),
						fmt.Sprintf("%v", task.Status),
					)

				}
			}
		}
		table.AddRow("", "", "", "")
	}

	return table
}
