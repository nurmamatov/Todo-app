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
}

func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefould("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(getOrReturnDefould("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToInt(getOrReturnDefould("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefould("POSTGRES_DATABASE", "assignes"))
	c.PostgresUser = cast.ToString(getOrReturnDefould("POSTGRES_USER", "khusniddin"))
	c.PostgresPassword = cast.ToString(getOrReturnDefould("POSTGRES_PASSWORD", "1234"))

	c.LogLevel = cast.ToString(getOrReturnDefould("LOG_LEVEL", "debug"))

	c.RPCPort = cast.ToString(getOrReturnDefould("RPC_PORT", ":50052"))

	return c
}

func getOrReturnDefould(key string, defouldValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defouldValue
}
