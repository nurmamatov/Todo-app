package main

import (
	"two_services/api/api"
	"two_services/api/config"
	"two_services/api/pkg/logger"
	"two_services/api/services"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

func main() {

	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	
	a := fileadapter.NewAdapter("./config/check.csv")
	
	casbinEnforcer, err := casbin.NewEnforcer(cfg.CasbinConfigPath, a)
	if err != nil {
		log.Error("new enforcer error", logger.Error(err))
		return
	}
	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		log.Error("casbin load policy error", logger.Error(err))
		return
	}
	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		ServiceManager: serviceManager,
		CasbinEnforcer: casbinEnforcer,
	})

	serviceManager.TaskService()
	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}

}
