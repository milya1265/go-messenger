package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrorTokenTimeOut = errors.New("manager: token timeout")
	ErrorTakeClaims   = errors.New("manager: error get user claims from token")
)

// TokenManager provides logic for JWT & Refresh tokens generation and parsing.
type TokenManager interface {
	ParseSub(accessToken string) (string, error)
	GenerateRefreshToken(userId string) (string, error)
	GenerateAccessToken(userId string) (string, error)
}

type Manager struct {
	signingKey      string
	ttlAccessToken  time.Duration
	ttlRefreshToken time.Duration
}

func NewManager(signingKey string, ttlRefresh time.Duration, ttlAccess time.Duration) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{
		signingKey:      signingKey,
		ttlAccessToken:  ttlAccess,
		ttlRefreshToken: ttlRefresh,
	}, nil
}

func (m *Manager) GenerateRefreshToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(m.ttlRefreshToken).Unix(),
		"sub": userId,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) GenerateAccessToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(m.ttlAccessToken).Unix(),
		"sub": userId,
	})
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) ParseSubject(Token string) (string, error) {
	token, err := jwt.Parse(Token, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrorTakeClaims
	}

	return claims["sub"].(string), nil
}

func (m *Manager) ParseExpiration(t string) (int64, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrorTakeClaims
	}

	return int64(claims["exp"].(float64)), nil
}
