package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment      string
	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string
	LogLevel         string
	RPCPort          string
	AssigneePort     int
	AssigneeHost     string
	EmailPort        int
	EmailHost        string
}

func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefould("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(getOrReturnDefould("POSTGRES_HOST", "database"))
	c.PostgresPort = cast.ToInt(getOrReturnDefould("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefould("POSTGRES_DATABASE", "database"))
	c.PostgresUser = cast.ToString(getOrReturnDefould("POSTGRES_USER", "khusniddin"))
	c.PostgresPassword = cast.ToString(getOrReturnDefould("POSTGRES_PASSWORD", "1234"))

	c.LogLevel = cast.ToString(getOrReturnDefould("LOG_LEVEL", "debug"))

	c.RPCPort = cast.ToString(getOrReturnDefould("RPC_PORT", ":50051"))
	c.AssigneeHost = cast.ToString(getOrReturnDefould("ASSIGNEE_SERVICE_HOST", "assignee-services"))
	c.AssigneePort = cast.ToInt(getOrReturnDefould("ASSIGNEE_SERVICE_PORT", 50052))

	c.EmailHost = cast.ToString(getOrReturnDefould("EMAIL_SERVICE_HOST", "email_service"))
	c.EmailPort = cast.ToInt(getOrReturnDefould("EMAIL_SERVICE_PORT", 9002))

	return c
}

func getOrReturnDefould(key string, defouldValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defouldValue
}
