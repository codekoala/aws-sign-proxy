package main

import "github.com/kelseyhightower/envconfig"

var config Config

type Config struct {
	Bind         string `default:":8080"`
	ExtraHeaders map[string]string
	TargetProto  string `default:"https"`
	TargetHost   string
	Region       string
	Provider     string
}

func init() {
	err := envconfig.Process("AWS_SIGN_PROXY", &config)
	if err != nil {
		panic(err)
	}
}
