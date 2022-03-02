package repo

import (
	pbe "two_services/task-services/genproto/email"
	pb "two_services/task-services/genproto/task"
)

type TaskStorageI interface {
	Create(*pb.CreateTaskReq) (*pb.TaskRes, error)
	Get(*pb.GetAndDeleteTask) (*pb.TaskRes, error)
	Update(*pb.UpdateTaskReq) (*pb.TaskRes, error)
	Delete(*pb.GetAndDeleteTask) (*pb.ErrOrStatus, error)
	List(*pb.LimAndPage) (*pb.TasksList, error)
	ListOverdue(*pb.Empty) (*pb.TasksList, error)
}

type EmailStorageI interface {
	Send(*pbe.Email) (*pbe.Empty, error)
}
