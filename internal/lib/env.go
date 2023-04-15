package lib

import (
	"os"

	"github.com/mitchellh/mapstructure"
)

type Env struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	PORT        string `mapstructure:"PORT"`

	AWS_REGION            string `mapstructure:"AWS_REGION"`
	AWS_ACCESS_KEY_ID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWS_SECRET_ACCESS_KEY string `mapstructure:"AWS_SECRET_ACCESS_KEY"`

	COGNITO_USER_POOL_ID string `mapstructure:"COGNITO_USER_POOL_ID"`
	COGNITO_CLIENT_ID    string `mapstructure:"COGNITO_CLIENT_ID"`

	ACCESS_KEY_SECRET    string `mapstructure:"ACCESS_KEY_SECRET"`
	REFRESH_TOKEN_SECRET string `mapstructure:"REFRESH_TOKEN_SECRET"`

	EMAIL_SERVICE_HOST string `mapstructure:"EMAIL_SERVICE_HOST"`
	APPLICATION_DOMAIN string `mapstructure:"APPLICATION_DOMAIN"`
}

var globalEnv = Env{}

func NewEnv(logger Logger) *Env {
	envVars := make(map[string]string)

	for _, key := range []string{
		"ENVIRONMENT",
		"PORT",
		"AWS_REGION",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"COGNITO_USER_POOL_ID",
		"COGNITO_CLIENT_ID",
		"ACCESS_KEY_SECRET",
		"REFRESH_TOKEN_SECRET",
		"EMAIL_SERVICE_HOST",
		"APPLICATION_DOMAIN",
	} {
		envVars[key] = os.Getenv(key)
	}

	err := mapstructure.Decode(envVars, &globalEnv)

	if err != nil {
		logger.Fatal("enable to map env variables", err)
	}
	return &globalEnv
}
