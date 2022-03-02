package main

import (
	"net"
	"two_services/assignee-services/config"
	"two_services/assignee-services/pkg/db"
	"two_services/assignee-services/pkg/logger"

	pb "two_services/assignee-services/genproto"
	"two_services/assignee-services/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "assignee-service")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	userService := service.NewUserService(connDB, log)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pb.RegisterUserServiceServer(s, userService)
	log.Info("main: server runnning",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

}
