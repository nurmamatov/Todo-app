package service

import (
	"context"

	// "fmt"
	pb "two_services/task-services/genproto/task"
	"two_services/task-services/genproto/user"

	pbe "two_services/task-services/genproto/email"
	l "two_services/task-services/pkg/logger"
	grpcClient "two_services/task-services/service/grpc_client"

	"two_services/task-services/storage"

	"github.com/jmoiron/sqlx"
)

type TaskService struct {
	storage storage.IStorage
	client  grpcClient.IService
	logger  l.Logger
}

func NewTaskService(db *sqlx.DB, log l.Logger, client grpcClient.IService) *TaskService {
	return &TaskService{
		storage: storage.NewStoragePg(db),
		client:  client,
		logger:  log,
	}
}

func (s *TaskService) Create(ctx context.Context, req *pb.CreateTaskReq) (*pb.TaskRes, error) {
	res, err := s.storage.Task().Create(req)
	if err != nil {
		s.logger.Error("error creat a task ", l.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *TaskService) Get(cft context.Context, req *pb.GetAndDeleteTask) (*pb.TaskRes, error) {
	res, err := s.storage.Task().Get(req)
	if err != nil {
		s.logger.Error("error get a task ", l.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *TaskService) Update(cft context.Context, req *pb.UpdateTaskReq) (*pb.TaskRes, error) {
	res, err := s.storage.Task().Update(req)
	if err != nil {
		s.logger.Error("error update a task ", l.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *TaskService) Delete(cft context.Context, req *pb.GetAndDeleteTask) (*pb.ErrOrStatus, error) {
	res, err := s.storage.Task().Delete(req)
	if err != nil {
		s.logger.Error("error delete a task ", l.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *TaskService) List(cft context.Context, req *pb.LimAndPage) (*pb.TasksList, error) {
	res, err := s.storage.Task().List(req)
	if err != nil {
		s.logger.Error("error get a list ", l.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *TaskService) ListOverdue(cft context.Context, req *pb.Empty) (*pb.TasksList, error) {
	res, err := s.storage.Task().ListOverdue(req)
	if err != nil {
		s.logger.Error("error get a listoverdue ", l.Error(err))
		return nil, err
	}
	
	var emails []string
	var phones []string
	for i, val := range res.Tasks {
		user, err := s.client.UserService().Get(cft, &user.GetOrDeleteUser{Id: val.AssigneeId})
		if err != nil {
			return nil, err
		}
		emails = append(emails, user.Email)
		phones = append(phones, user.PhoneNum[i].Phone)
	}
	email := &pbe.Email{
		Id:         "",
		Body:       "Your task has passed the deadline",
		Phone:      phones,
		Subject: "Warning",
		Recipients: emails,
	}
	_, err = s.client.EmailService().SendEmail(cft, email)
	if err != nil {
		return nil, err
	}
	for _,j := range phones {

		sms := &pbe.Sms{
			Id:         "",
			Body:       "Your task has passed the deadline",
			Phone:      j,
		}
		_, err = s.client.EmailService().SendSms(cft, sms)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
