package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"

	"github.com/codekoala/aws-sign-proxy"
)

var config aws_sign_proxy.Config

func init() {
	err := envconfig.Process("AWS_SIGN_PROXY", &config)
	if err != nil {
		panic(err)
	}

	if config.Region == "" {
		config.Region = os.Getenv("AWS_DEFAULT_REGION")
	}
}
