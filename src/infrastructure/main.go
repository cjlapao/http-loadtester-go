package infrastructure

import (
	"github.com/cjlapao/common-go/log"
)

var logger = log.Get()

func Init() {
	NewServiceProvider()
}
