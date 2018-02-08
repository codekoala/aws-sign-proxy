package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/codekoala/aws-sign-proxy"
)

var log *zap.Logger

func init() {
	var err error

	if log, err = zap.NewProduction(); err != nil {
		panic(err)
	}
}

func main() {
	log.Info("getting access key ID and secret from environment")
	creds := credentials.NewEnvCredentials()
	signer := v4.NewSigner(creds)
	req_signer := aws_sign_proxy.NewRequestSigner(log, config, signer)

	http.Handle(config.MetricsEndpoint, promhttp.Handler())
	http.HandleFunc(config.HealthzEndpoint, healthz)
	http.HandleFunc("/", req_signer.Proxy)

	svr := &http.Server{Addr: config.Bind}
	go serve(svr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutdown signal received, exiting...")
	svr.Shutdown(context.Background())
}

// healthz is just a simple health check for the service
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// serve begins accepting and handling requests
func serve(svr *http.Server) {
	log.Info("accepting connections", zap.String("addr", config.Bind))
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal("error serving requests", zap.Error(err))
	}
}
