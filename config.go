package aws_sign_proxy

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

var config Config

type Config struct {
	Bind         string `default:":8080"`
	ExtraHeaders map[string]string
	BlockHeaders []string
	TargetProto  string `default:"https"`
	TargetHost   string
	Region       string
	Provider     string

	HealthzEndpoint string `default:"/_healthz"`
	MetricsEndpoint string `default:"/_metrics"`
}

func init() {
	err := envconfig.Process("AWS_SIGN_PROXY", &config)
	if err != nil {
		panic(err)
	}

	if config.Region == "" {
		config.Region = os.Getenv("AWS_DEFAULT_REGION")
	}
}
