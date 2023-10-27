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
	SaveReadStatus(message *ReadMessage) error
	GetChatMessages(chatID, limit, offset int, userID string) ([]*Message, error)
	DeleteMessage(messageID int, userID string) error
	EditMessage(messageID int, text, userID string) error
}

func (s *service) NewMessage(message *Message) (*Message, error) {
	c, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if len(message.Text) > 512 {
		s.logger.Error(errors.New("message too long, max length = 512"))
		return nil, errors.New("message too long, max length = 512")
	}

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

	if len(members) == 0 {
		s.logger.Error(errors.New("count of members can't be zero"))
		return nil, errors.New("count of members can't be zero")
	}

	if chat.IsDirect == false && chat.Title == "" {
		s.logger.Error(errors.New("not direct chat must have title"))
		return nil, errors.New("not direct chat must have title")
	}

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

	s.logger.Debug(res)

	err = s.repo.NewMembers(c, res.Id, members)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	s.logger.Debug(members)

	res.Members = members
	res.IsDirect = chat.IsDirect

	return res, nil
}

func (s *service) GetChatMembers(chatID int) ([]string, error) {
	c, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.repo.GetChatMembers(c, chatID)
}

func (s *service) SaveReadStatus(message *ReadMessage) error {
	c, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	err := s.repo.CheckReadStatus(c, message)

	if err != nil {
		if errors.Is(err, errors.New("read status not found in this chat")) {
			return s.repo.StoreReadStatus(c, message)
		}
		s.logger.Error(err)
		return err
	}

	return s.repo.UpdateReadStatus(c, message)
}

func (s *service) GetChatMessages(chatID, limit, offset int, userID string) ([]*Message, error) {
	c, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	err := s.repo.CheckChatMember(c, chatID, userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	messages, err := s.repo.GetChatMessages(c, chatID, limit, offset)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return messages, nil
}

func (s *service) EditMessage(messageID int, text, userID string) error {
	c, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	message, err := s.repo.GetMessageByID(c, messageID)
	if err != nil {
		return err
	}

	err = s.repo.CheckAuthorMessage(c, message.Id, userID)
	if err != nil {
		return err
	}

	return s.repo.EditTextMessage(c, message.Id, text)
}

func (s *service) DeleteMessage(messageID int, userID string) error {
	c, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	message, err := s.repo.GetMessageByID(c, messageID)
	if err != nil {
		return err
	}

	err = s.repo.CheckChatMember(c, message.ChatID, userID)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return s.repo.DeleteMessage(c, message.Id)
}
