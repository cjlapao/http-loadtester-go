package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/version"

	"github.com/cjlapao/common-go/log"
)

var logger = log.Get()
var versionSvc = version.Get()

func main() {
	versionSvc.Minor = 1
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

	file := helper.GetFlagValue("file", "")

	if file != "" {
		err := ExecuteFromFile(file)
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
