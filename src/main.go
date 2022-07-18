package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/restapi"
	"github.com/cjlapao/common-go/version"
	"github.com/cjlapao/http-loadtester-go/controllers"
	"github.com/cjlapao/http-loadtester-go/database"
	"github.com/cjlapao/http-loadtester-go/infrastructure"
	"github.com/cjlapao/http-loadtester-go/usecases"
)

var ver = "0.0.30"
var logger = log.Get()
var versionSvc = version.Get()

func main() {
	logger.WithTimestamp()
	SetVersion()
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
		Init()
	} else {
		if file != "" {
			infrastructure.Init()
			err := usecases.ExecuteFromFile(file)
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

func Init() {
	listener := restapi.GetHttpListener()
	listener.AddLogger().AddHealthCheck()
	listener.AddController(controllers.LoadController, "/start", "POST")
	listener.AddController(controllers.StartLoadFileController, "/start/file", "POST")
	listener.AddController(controllers.StartLoadMarkdownController, "/start/markdown", "POST")

	database.Init()
	infrastructure.Init()
	listener.Start()
}

func SetVersion() {
	versionSvc.Name = "HTTP LoadTester"
	versionSvc.Author = "carlos Lapao"
	versionSvc.License = "MIT"
	strVer, err := version.FromString(ver)
	if err == nil {
		versionSvc.Major = strVer.Major
		versionSvc.Minor = strVer.Minor
		versionSvc.Build = strVer.Build
		versionSvc.Rev = strVer.Rev
	}
}
