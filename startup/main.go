package startup

import (
	"github.com/cjlapao/common-go/restapi"
	"github.com/cjlapao/http-loadtester-go/controller"
)

func Init() {
	listener := restapi.GetHttpListener()
	listener.AddLogger().AddHealthCheck()
	listener.AddController(controller.LoadController, "/start", "POST")
	listener.AddController(controller.StartLoadFileController, "/start/file", "POST")
	listener.AddController(controller.StartLoadMarkdownController, "/start/markdown", "POST")
	listener.Start()
}
