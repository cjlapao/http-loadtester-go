package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/restapi"
	"github.com/cjlapao/common-go/version"
	"github.com/cjlapao/http-loadtester-go/jobs"
	"github.com/cjlapao/http-loadtester-go/startup"
)

var logger = log.Get()
var versionSvc = version.Get()
var api = restapi.NewHttpListener()

func main() {
	logger.WithTimestamp()
	versionSvc.Name = "HTTP LoadTester"
	versionSvc.Author = "carlos Lapao"
	versionSvc.License = "MIT"
	versionSvc.Minor = 2
	versionSvc.Build = 0
	getVersion := helper.GetFlagSwitch("version", false)
	if getVersion {
		format := helper.GetFlagValue("o", "json")
		switch strings.ToLower(format) {
		case "json":
			fmt.Println(versionSvc.PrintVersion(int(version.JSON)))
		case "yaml":
			fmt.Println(versionSvc.PrintVersion(int(version.JSON)))
		default:
			fmt.Println("Please choose a valid format, this can be either json or yaml")
		}
		os.Exit(0)
	}
	versionSvc.PrintAnsiHeader()

	file := ""
	if helper.FileExists("config.yml") {
		file = "config.yml"
	} else {
		file = helper.GetFlagValue("file", "")
	}

	apiMode := helper.GetFlagSwitch("api", false)

	if apiMode {
		startup.Init()
	} else {
		if file != "" {
			err := jobs.ExecuteFromFile(file)
			if err != nil {
				logger.Error("There was an error processing the file")
				os.Exit(1)
			}
			logger.Success("Finished, bye!!!")
			os.Exit(0)
		}

		url := helper.GetFlagValue("target", "")
		if url == "" {
			logger.Error("Missing url to target")
			os.Exit(1)
		}

		fmt.Print("Finished")
	}
}
