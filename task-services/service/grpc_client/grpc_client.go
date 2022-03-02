package grpcclient

import (
	"fmt"
	"two_services/task-services/config"
	pb "two_services/task-services/genproto/user"
	pbe "two_services/task-services/genproto/email"

	"google.golang.org/grpc"
)

type IService interface {
	UserService() pb.UserServiceClient
	EmailService() pbe.EmailServiceClient
}

type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

func New(cfg config.Config) (*GrpcClient, error) {
	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.AssigneeHost, cfg.AssigneePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	
	connEmail, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.EmailHost, cfg.EmailPort),
		grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"user_service": pb.NewUserServiceClient(connUser),
			"email_service": pbe.NewEmailServiceClient(connEmail),
		},
	}, nil
}


func (g *GrpcClient) UserService() pb.UserServiceClient {
	return g.connections["user_service"].(pb.UserServiceClient)
}

func (g *GrpcClient) EmailService() pbe.EmailServiceClient {
	return g.connections["email_service"].(pbe.EmailServiceClient)
}