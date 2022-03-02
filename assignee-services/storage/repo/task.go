package repo

import (
	pb "two_services/assignee-services/genproto"
)

type UserStorageI interface {
	Create(*pb.CreateUserReqWithCode) (*pb.UserRes, error)
	Get(*pb.GetOrDeleteUser) (*pb.UserRes, error)
	Update(*pb.UpdateUserReq) (*pb.UserRes, error)
	Delete(*pb.GetOrDeleteUser) (*pb.ErrOrStatus, error)
	List(*pb.Empty) (*pb.UsersList, error)
	ListEmail(*pb.Ids) (*pb.Emails, error)
	ChekUser(*pb.EmailWithUsername) (*pb.Bool, bool)
	GetEmail(req *pb.Email) (string, error)
	Regist(req *pb.CreateUserReqWithCode) (*pb.Mess, error)
	Verfy(req *pb.Check) (*pb.CreateUserReq, error)
	Login(req *pb.EmailWithPassword) (*pb.UserRes, error)
	UpdateToken(req *pb.TokensReq) (*pb.Tokens, error)
	Filtr(req *pb.FiltrReq) (*pb.UsersList, error)
}
