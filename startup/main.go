package startup

import (
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/restapi"
	"github.com/cjlapao/http-loadtester-go/controller"
	"github.com/cjlapao/http-loadtester-go/dbctx"
)

var contextSvc = execution_context.Get()

func Init() {
	listener := restapi.GetHttpListener()
	listener.AddLogger().AddHealthCheck()
	listener.AddController(controller.LoadController, "/start", "POST")
	listener.AddController(controller.StartLoadFileController, "/start/file", "POST")
	listener.AddController(controller.StartLoadMarkdownController, "/start/markdown", "POST")

	dbctx.Init()
	listener.Start()
}
