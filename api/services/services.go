package services

import (
	"fmt"
	"two_services/api/config"
	pbemail "two_services/api/genproto/email"
	pbtask "two_services/api/genproto/task"
	pbuser "two_services/api/genproto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	TaskService() pbtask.TaskServiceClient
	UserService() pbuser.UserServiceClient
	EmailService() pbemail.EmailServiceClient
}

type serviceManager struct {
	taskService  pbtask.TaskServiceClient
	userService  pbuser.UserServiceClient
	emailService pbemail.EmailServiceClient
}

func (s *serviceManager) TaskService() pbtask.TaskServiceClient {
	return s.taskService
}
func (s *serviceManager) UserService() pbuser.UserServiceClient {
	return s.userService
}
func (s *serviceManager) EmailService() pbemail.EmailServiceClient {
	return s.emailService
}

func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	connTask, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.TaskServiceHost, conf.TaskServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.UserServiceHost, conf.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	connEmail, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.EmailServiceHost, conf.EmailServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	serviceManager := &serviceManager{
		taskService: pbtask.NewTaskServiceClient(connTask),
		userService: pbuser.NewUserServiceClient(connUser),
		emailService: pbemail.NewEmailServiceClient(connEmail),
	}

	return serviceManager, nil
}
