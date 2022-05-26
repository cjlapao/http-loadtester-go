package startup

import (
	"github.com/cjlapao/common-go/restapi"
	"github.com/cjlapao/http-loadtester-go/controller"
)

func Init() {
	listener := restapi.GetHttpListener()
	listener.AddJsonContent().AddLogger().AddHealthCheck()
	listener.AddController(controller.LoadController, "/load-tester/start", "POST")
	listener.Start()
}
