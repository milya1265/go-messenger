package messenger

import (
	"WSChats/pkg/logger"
	"context"
	"database/sql"
	"errors"
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
	NewChat(chat *NewChatReq) (*NewChatRes, error)
	GetChatMembers(chatID int) ([]string, error)
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

func (s *service) NewChat(chat *NewChatReq) (*NewChatRes, error) {
	c, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	members := chat.Members

	if chat.IsDirect == true && len(members) != 2 {
		s.logger.Error(errors.New("direct members must be two"))
		return nil, errors.New("direct members must be two")
	}

	if chat.IsDirect == true {
		chatID, err := s.repo.SearchDirectChat(c, members[0], members[1])
		if !errors.Is(sql.ErrNoRows, err) && err != nil {
			s.logger.Error(err)
			return nil, err
		}
		if chatID != 0 {
			return nil, errors.New("chat has already created")
		}
		chat.Title = ""
	}

	res, err := s.repo.NewChat(c, chat)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	err = s.repo.NewMembers(c, res.Id, members)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	res.Members = members
	res.IsDirect = chat.IsDirect

	return res, nil
}

func (s *service) GetChatMembers(chatID int) ([]string, error) {
	c, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.repo.GetChatMembers(c, chatID)
}
