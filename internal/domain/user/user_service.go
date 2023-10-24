package user

import (
	"WSChats/internal/adapters/api/DTO"
	"WSChats/internal/util"
	"WSChats/pkg/logger"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

var JWTkey = []byte("qwerty")

type service struct {
	repo    Repository
	timeout time.Duration
	logger  logger.Logger
}

func NewService(r *Repository, l *logger.Logger) Service {
	return &service{
		repo:    *r,
		timeout: time.Duration(2 * time.Second),
		logger:  *l,
	}
}

type Service interface {
	CreateUser(ctx context.Context, u *DTO.CreateUserReq) (*DTO.CreateUserRes, error)
	Login(ctx context.Context, u *DTO.GetUserByEmailReq) (*DTO.GetUserByEmailRes, error)
	//UpdateUser() (*User, error)
	//DeleteUser() error
}

func (s *service) CreateUser(ctx context.Context, req *DTO.CreateUserReq) (*DTO.CreateUserRes, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var u User
	UUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	u.UUID = UUID.String()
	u.Email = req.Email
	u.Username = req.Username

	hashedPass, err := util.HashPassword([]byte(req.Password))
	if err != nil {
		return nil, err
	}

	u.Password = string(hashedPass)

	res := &DTO.CreateUserRes{}

	res, err = s.repo.CreateUser(c, &u)
	if err != nil {
		return nil, err
	}

	return res, nil
}

var ErrPasswordNotCompare = errors.New("password is not compare")

func (s *service) Login(ctx context.Context, req *DTO.GetUserByEmailReq) (*DTO.GetUserByEmailRes, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	userDb, err := s.repo.GetUserByEmail(c, req)
	if err != nil {
		return nil, err
	}

	if err := util.ComparePassword([]byte(req.Password), []byte(userDb.Password)); err != nil {
		return nil, err
	}

	var res = &DTO.GetUserByEmailRes{
		UUID:     userDb.UUID,
		Email:    userDb.Email,
		Username: userDb.Email,
	}
	s.logger.Debug("В сервисе всё ок")

	return res, nil
}
func (s *service) UpdateUser() {

}

func (s *service) DeleteUser() {

}
