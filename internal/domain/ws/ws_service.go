package ws

import (
	"WSChats/pkg/logger"
	"context"
	"time"
)

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
	NewMessage(message *Message) (*Message, error)
	GetUsernameByUUID(uuid string) (string, error)
}

func (s *service) NewMessage(message *Message) (*Message, error) {
	c, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.repo.NewMessage(c, message)
}

func (s *service) GetUsernameByUUID(uuid string) (string, error) {
	c, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.repo.GetUsernameByUUID(c, uuid)
}
