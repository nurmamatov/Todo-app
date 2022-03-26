package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Enivorentment    string
	TaskServiceHost  string
	UserServiceHost  string
	TaskServicePort  int
	UserServicePort  int
	EmailServiceHost string
	EmailServicePort int

	CtxTimeout int

	LogLevel string
	HTTPPort string

	SigninKey string
	RedisHost string
	RedisPort int

	CasbinConfigPath string
}

func Load() Config {
	c := Config{}

	c.Enivorentment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8080"))
	c.TaskServiceHost = cast.ToString(getOrReturnDefault("TASK_SERVICE_HOST", "task-services"))
	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "assignee-services"))
	c.TaskServicePort = cast.ToInt(getOrReturnDefault("TASK_SERVICE_PORT", 50051))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 50052))
	c.EmailServiceHost = cast.ToString(getOrReturnDefault("EMAIL_SERVICE_HOST", "email_service"))
	c.EmailServicePort = cast.ToInt(getOrReturnDefault("EMAIL_SERVICE_PORT", 9002))

	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))
	c.CasbinConfigPath = cast.ToString(getOrReturnDefault("CASBIN_CONFIG_PATH", "config/rbac_model.conf"))
	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
