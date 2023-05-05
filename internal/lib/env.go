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

	S3_BUCKET_NAME string `mapstructure:"S3_BUCKET_NAME"`
	DATABASE_DSN   string `mapstructure:"DATABASE_DSN"`

	REDIS_HOST string `mapstructure:"REDIS_HOST"`
	REDIS_PORT string `mapstructure:"REDIS_PORT"`

	AUDITOR_INSTALL_NAME string `mapstructure:"AUDITOR_INSTALL_NAME"`
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
		"S3_BUCKET_NAME",
		"DATABASE_DSN",
		"REDIS_HOST",
		"REDIS_PORT",
		"AUDITOR_INSTALL_NAME",
	} {

		envVars[key] = os.Getenv(key)
	}

	err := mapstructure.Decode(envVars, &globalEnv)

	if err != nil {
		logger.Fatal("enable to map env variables", err)
	}
	return &globalEnv
}
