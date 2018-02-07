package main

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"go.uber.org/zap"

	"github.com/codekoala/aws-sign-proxy"
)

var log *zap.Logger

func main() {
	var err error

	if log, err = zap.NewProduction(); err != nil {
		panic(err)
	}

	log.Info("getting access key ID and secret from environment")
	creds := credentials.NewEnvCredentials()
	signer := v4.NewSigner(creds)

	http.HandleFunc("/", aws_sign_proxy.SignRequest(log, config, signer))

	log.Info("accepting connections", zap.String("addr", config.Bind))
	if err = http.ListenAndServe(config.Bind, nil); err != nil {
		log.Fatal("error serving requests", zap.Error(err))
	}
}
