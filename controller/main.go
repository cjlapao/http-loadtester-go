package controller

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"net/http"

	commonLogger "github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/version"
	"github.com/cjlapao/http-loadtester-go/startup"
	"github.com/gorilla/mux"
)

var logger = commonLogger.Get()
var versionSvc = version.Get()

// var router mux.Router
var serviceProvider = startup.CreateProvider()

// Controller Controller structure
type Controller struct {
	Router *mux.Router
}

var globalController *Controller

// NewAPIController  Creates a new controller
func NewAPIController(router *mux.Router) *Controller {
	if globalController != nil {
		return globalController
	}

	controller := Controller{
		Router: router,
	}

	controller.Router.HandleFunc(serviceProvider.Context.ApiPrefix+"/login", controller.Login).Methods("POST")
	controller.Router.HandleFunc(serviceProvider.Context.ApiPrefix+"/validate", controller.Validate).Methods("GET")

	globalController = &controller
	return globalController
}

func RestApiModuleProcessor() {
	logger.Notice("Starting Go Rest API v%v", versionSvc.String())
	if serviceProvider.Context.BackendEnabled {
		logger.Info("Found MongoDB connection, enabling MongoDb backend...")

	}
	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Homepage!")
	fmt.Println("endpoint Hit: homepage")
}

func handleRequests() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)
	router.HandleFunc(serviceProvider.Context.ApiPrefix+"/", homePage)
	_ = NewAPIController(router)

	// Creating and starting the http server
	srv := &http.Server{
		Addr:    ":" + serviceProvider.Context.ApiPort,
		Handler: router,
	}

	go func() {
		logger.Info("Api listening on http://localhost:" + serviceProvider.Context.ApiPort + serviceProvider.Context.ApiPrefix)
		logger.Success("Finished Initiating http server")
		if err := srv.ListenAndServe(); err != nil {
			logger.FatalError(err, "There was an error")
		}
	}()

	if serviceProvider.Context.TLS {
		cert, err := tls.X509KeyPair([]byte(serviceProvider.Context.TLSCertificate), []byte(serviceProvider.Context.TLSPrivateKey))
		if err == nil {
			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{cert},
			}

			sslSrv := &http.Server{
				Addr:      ":" + serviceProvider.Context.TLSApiPort,
				TLSConfig: tlsConfig,
				Handler:   router,
			}

			go func() {
				logger.Info("Api listening on https://localhost:" + serviceProvider.Context.TLSApiPort + serviceProvider.Context.ApiPrefix)
				logger.Success("Finished Initiating https server")
				if err := sslSrv.ListenAndServeTLS("", ""); err != nil {
					logger.FatalError(err, "There was an error")
				}
			}()
		} else {
			logger.Error("There was an error reading the certificates to enable HTTPS")
		}
	}

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	logger.Info("Server shutdown")
	// http.ListenAndServeTLS(":10001", "./ssl/local-cluster.internal.crt", "./ssl/local-cluster.internal.key", router)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
