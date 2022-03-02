package v1

import (
	"two_services/api/config"
	"two_services/api/pkg/logger"
	"two_services/api/services"
	"two_services/api/api/token"
	// "two_services/api/storage/repo"
)

type HandlerV1 struct {
	log             logger.Logger
	serviceManager  services.IServiceManager
	cfg             config.Config
	// inMemoryStorage repo.InMemoryStorageI
	jwtHandler      token.JWTHandler
}

type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	JwtHandler token.JWTHandler
}

func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
	}
}
