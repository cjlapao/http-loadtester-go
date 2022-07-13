package main

import "github.com/cjlapao/common-go/version"

var ver = "0.2.0.0"

func SetVersion() {
	versionSvc.Name = "HTTP LoadTester"
	versionSvc.Author = "carlos Lapao"
	versionSvc.License = "MIT"
	strVer, err := version.FromString(ver)
	if err == nil {
		versionSvc.Major = strVer.Major
		versionSvc.Minor = strVer.Minor
		versionSvc.Build = strVer.Build
		versionSvc.Build = strVer.Build
	}
}
