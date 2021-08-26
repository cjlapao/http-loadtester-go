package jobs

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cjlapao/common-go/helper"

	"github.com/cjlapao/markdown-go"
	"gopkg.in/yaml.v2"
)

// ExportReportToFile Exports a Job Report to a File
func (j *JobOperation) ExportReportToFile() {
	content := j.MarkDown()
	name := j.Name
	if *name != j.ID {
		*name += "_" + j.ID
	}
	helper.WriteToFile(content, "job_"+*name+"_report.md")
}

// ExportOutputToFile Exports a job result to a file
func (j *JobOperation) ExportOutputToFile() {
	stringContent, err := yaml.Marshal(j.Result)
	if err != nil {
		return
	}
	name := j.Name
	if *name != j.ID {
		*name += "_" + j.ID
	}
	helper.WriteToFile(string(stringContent), "job_"+*name+"_output.yml")
}

// MarkDown Generates a Job Markdown report
func (j *JobOperation) MarkDown() string {
	md := markdown.CreateDocument()
	header := md.CreateHeader()
	header.H1(fmt.Sprintf("HTTP Load Tester Report for %v", j.Target.URL))

	htb := md.CreateTextBlock()
	htb.AddLine(fmt.Sprintf("Job Type: %v", j.Type))
	htb.AddLine(fmt.Sprintf("Block Type: %v", j.OperationType))
	htb.AddLine(fmt.Sprintf("Task Type: %v", j.Options.BlockType))
	htb.AddLine(fmt.Sprintf("Timeout: %v", fmt.Sprint(time.Duration(j.Options.Timeout)*time.Second)))
	htb.AddLine(fmt.Sprintf("Method: %v", j.Target.Method))
	if j.Target.Body != "" {
		htb.AddLine("Contains Body: Yes")
	} else {
		htb.AddLine("Contains Body: No")
	}
	htb.AddLine(fmt.Sprintf("Duration: %.2f seconds", j.Result.TotalDurationInSeconds))
	htb.AddLine(fmt.Sprintf("Total Calls: %v, Succeeded: %v, Failed: %v", j.Result.TotalCalls, j.Result.TotalSucceededCalls, j.Result.TotalFailedCalls))
	htb.AddLine(fmt.Sprintf("Average Block Duration: %.4f seconds", j.Result.AverageBlockDuration))
	htb.AddLine(fmt.Sprintf("Average Call Duration: %.4f seconds", j.Result.AverageCallDuration))
	htb.AddLine(fmt.Sprintf("Authentication: %v", strconv.FormatBool(j.Authenticated())))
	htb.AddLine(fmt.Sprintf("Time Taken: %.2f seconds", j.Result.TimeTaken.Seconds()))
	htb.NewLine().NewLine()
	statusDetailsHeader := md.CreateHeader()
	statusDetailsHeader.H2("Status Results Details")
	j.generateTaskTable(md)
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
