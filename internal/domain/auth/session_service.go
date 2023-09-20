package auth

import (
	"WSChats/pkg/logger"
	"context"
	"time"
)

var JWTkey = []byte("qwerty")

type service struct {
	repo       Repository
	timeout    time.Duration
	logger     logger.Logger
	managerJWT Manager
}

func NewService(r *Repository, l *logger.Logger, m *Manager) Service {
	return &service{
		repo:       *r,
		timeout:    time.Duration(2 * time.Second),
		logger:     *l,
		managerJWT: *m,
	}
}

type Service interface {
	NewSession(ctx context.Context, uuid string) (string, error)
	Authorize(ctx context.Context, access string) (string, string, error)
}

func (s *service) NewSession(ctx context.Context, uuid string) (string, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	newAccessToken, err := s.managerJWT.GenerateAccessToken(uuid)
	if err != nil {
		return "", err
	}
	newRefreshToken, err := s.managerJWT.GenerateRefreshToken(uuid)
	if err != nil {
		return "", err
	}

	sess, err := s.repo.NewAccessAndRefreshToken(c, uuid, newAccessToken, newRefreshToken)
	if err != nil {
		return "", err
	}

	return sess.accessToken, nil
}

func (s *service) Authorize(ctx context.Context, access string) (string, string, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	uuid, err := s.managerJWT.ParseSubject(access)
	if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}
	sess, err := s.repo.GetSessionByID(c, uuid)
	if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}

	exp, err := s.managerJWT.ParseExpiration(sess.accessToken)
	if err != nil {
		s.logger.Error(err.Error())
		return "", "", err
	}
	if exp < time.Now().Unix() {
		exp, err := s.managerJWT.ParseExpiration(sess.refreshToken)
		if err != nil {
			s.logger.Error(err.Error())
			return "", "", err
		}
		if exp < time.Now().Unix() {
			if err != nil {
				s.logger.Error(err.Error())
				return "", "", ErrorTokenTimeOut
			}
		} else {
			newAccess, err := s.managerJWT.GenerateAccessToken(uuid)
			if err != nil {
				s.logger.Error(err.Error())
				return "", "", err
			}
			sess, err := s.repo.NewAccessToken(c, uuid, newAccess)
			if err != nil {
				s.logger.Error(err.Error())
				return "", "", err
			}

			return sess.uuid, sess.accessToken, nil
		}
	}

	return uuid, access, nil
}

//newSess, err := s.repo.NewAccessAndRefreshToken(c, uuid, sess.accessToken, sess.refreshToken)
//if err != nil{
//return "", "", err
//}
//return newSess.uuid, newSess.accessToken, nil
