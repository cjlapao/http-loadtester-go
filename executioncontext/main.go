package executioncontext

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/log"
)

var globalContext *Context
var logger = log.Get()

// Context entity
type Context struct {
	MongoConnectionString string `json:"mongoConnectionString"`
	Database              string `json:"database"`
	ShowHelp              bool   `json:"help"`
	BackendEnabled        bool   `json:"backendEnabled"`
	ApiPrefix             string `json:"apiPrefix"`
	ApiPort               string `json:"apiPort"`
	TLS                   bool   `json:"tls"`
	TLSApiPort            string `json:"tlsApiPort"`
	TLSCertificate        string `json:"tlsCertificate"`
	TLSPrivateKey         string `json:"tlsPrivateKey"`
}

func Get() *Context {
	if globalContext != nil {
		return globalContext
	}

	logger.Debug("Creating Execution Context")
	globalContext = &Context{
		ShowHelp: helper.GetFlagSwitch("help", false),
	}

	globalContext.Getenv()

	return globalContext
}

// Getenv gets the environment variables for the entities
func (e *Context) Getenv() {

	e.Database = os.Getenv("RESTAPI_DATABASENAME")
	e.MongoConnectionString = os.Getenv("RESTAPI_MONGO_CONNECTION_STRING")
	e.ApiPrefix = os.Getenv("RESTAPI_API_PREFIX")
	e.ApiPort = os.Getenv("RESTAPI_API_PORT")
	e.TLS = false
	e.TLSApiPort = os.Getenv("RESTAPI_TLS_API_PORT")

	if os.Getenv("RESTAPI_TLS") != "" && strings.ToLower(os.Getenv("RESTAPI_TLS")) == "true" {
		e.TLS = true
	}

	if os.Getenv("RESTAPI_TLS_CERTIFICATE") != "" {
		decodedCert, _ := base64.StdEncoding.DecodeString(os.Getenv("RESTAPI_TLS_CERTIFICATE"))
		e.TLSCertificate = string(decodedCert)
	}

	if os.Getenv("RESTAPI_TLS_PRIVATE_KEY") != "" {
		decodedPrivateKey, _ := base64.StdEncoding.DecodeString(os.Getenv("RESTAPI_TLS_PRIVATE_KEY"))
		e.TLSPrivateKey = string(decodedPrivateKey)
	}

	if e.MongoConnectionString == "" {
		e.BackendEnabled = false
	} else {
		e.BackendEnabled = true
	}

	if e.Database == "" {
		e.Database = "restapi_demo_db"
	}

	if e.ApiPort == "" {
		e.ApiPort = "5000"
	}

	if e.TLSApiPort == "" {
		e.TLSApiPort = "5001"
	}

	// if the tls is enabled but no certificate or private key then we will disable it
	if e.TLSCertificate == "" || e.TLSPrivateKey == "" {
		e.TLS = false
	}

	e.ApiPrefix = strings.TrimSuffix(e.ApiPrefix, "/")
}
