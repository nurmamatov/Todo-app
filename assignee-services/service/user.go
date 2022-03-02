package service

import (
	"context"
	"errors"
	pb "two_services/assignee-services/genproto"
	l "two_services/assignee-services/pkg/logger"

	"two_services/assignee-services/storage"

	"github.com/jmoiron/sqlx"
)

type UserService struct {
	storage storage.IStorage
	logger  l.Logger
}

func NewUserService(db *sqlx.DB, log l.Logger) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}
func (s *UserService) Verfy(ctx context.Context, req *pb.Check) (*pb.CreateUserReq, error) {
	res, err := s.storage.User().Verfy(req)
	if err != nil {
		return nil, err
	}
	return res, err
}
func (s *UserService) Registr(ctx context.Context, req *pb.CreateUserReqWithCode) (*pb.Mess, error) {
	emailwithusername := pb.EmailWithUsername{Email: req.Email, Username: req.Username}
	_, err := s.ChekUser(ctx, &emailwithusername)
	if err != nil {
		return nil, err
	}
	res, err := s.storage.User().Regist(req)
	if err != nil {
		s.logger.Error("error registr a user ", l.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *UserService) Get(cft context.Context, req *pb.GetOrDeleteUser) (*pb.UserRes, error) {
	res, err := s.storage.User().Get(req)
	if err != nil {
		s.logger.Error("error get a user ", l.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *UserService) Update(cft context.Context, req *pb.UpdateUserReq) (*pb.UserRes, error) {
	emailwithusername := pb.EmailWithUsername{Email: req.Email, Username: req.Username}
	_, err := s.ChekUser(cft, &emailwithusername)
	if err != nil {
		return nil, err
	}
	res, err := s.storage.User().Update(req)
	if err != nil {
		s.logger.Error("error update a user ", l.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *UserService) Delete(cft context.Context, req *pb.GetOrDeleteUser) (*pb.ErrOrStatus, error) {
	res, err := s.storage.User().Delete(req)
	if err != nil {
		s.logger.Error("error delete a user ", l.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *UserService) List(cft context.Context, req *pb.Empty) (*pb.UsersList, error) {
	res, err := s.storage.User().List(req)
	if err != nil {
		s.logger.Error("error get a list ", l.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *UserService) ListEmail(cft context.Context, req *pb.Ids) (*pb.Emails, error) {
	res, err := s.storage.User().ListEmail(req)
	if err != nil {
		s.logger.Error("error get a list emails ", l.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *UserService) ChekUser(cft context.Context, req *pb.EmailWithUsername) (*pb.Bool, error) {
	res, bul := s.storage.User().ChekUser(req)
	if bul == bool(true) {
		return nil, errors.New("while Chek error")
	} else {
		return res, nil
	}
}
func (s *UserService) Login(cft context.Context, req *pb.EmailWithPassword) (*pb.UserRes, error) {
	res, err := s.storage.User().Login(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (s *UserService) Create(cft context.Context, req *pb.CreateUserReqWithCode) (*pb.UserRes, error) {
	res, err := s.storage.User().Create(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *UserService) UpdateToken(cft context.Context, req *pb.TokensReq) (*pb.Tokens, error) {
	res, err := s.storage.User().UpdateToken(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (s *UserService) Filtr(cft context.Context, req *pb.FiltrReq) (*pb.UsersList, error) {
	res,err := s.storage.User().Filtr(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
