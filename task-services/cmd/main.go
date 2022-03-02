package main

import (
	"net"
	"two_services/task-services/config"
	"two_services/task-services/pkg/db"
	"two_services/task-services/pkg/logger"
	grpcclient "two_services/task-services/service/grpc_client"

	pb "two_services/task-services/genproto/task"
	"two_services/task-services/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "task-service")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	grpcClient,err := grpcclient.New(cfg)
	if err!=nil {
		log.Fatal("grpc dial error",logger.Error(err))
	}

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	taskService := service.NewTaskService(connDB, log,grpcClient)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pb.RegisterTaskServiceServer(s, taskService)
	log.Info("main: server runnning",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

}
