package ws

import (
	"WSChats/pkg/logger"
	"context"
	"database/sql"
)

type repository struct {
	DB     *sql.DB
	logger *logger.Logger
}

func NewRepository(db *sql.DB, l *logger.Logger) Repository {
	return &repository{
		DB:     db,
		logger: l,
	}
}

type Repository interface {
	NewMessage(ctx context.Context, message *Message) (*Message, error)
	GetUsernameByUUID(ctx context.Context, uuid string) (string, error)
	//GetSessionByID(ctx context.Context, idUser string) (*session, error)
	//NewAccessAndRefreshToken(ctx context.Context, idUser string, access string, refresh string) (*session, error)
}

func (r *repository) NewMessage(ctx context.Context, message *Message) (*Message, error) {
	query := "INSERT INTO messages (sender, receiver, text, time) VALUES ($1, $2, $3, $4) RETURNING id;" //RETURNING

	row := r.DB.QueryRowContext(ctx, query, message.Sender, message.Receiver, message.Text, message.Time)

	var idMes int

	if err := row.Scan(&idMes); err != nil {
		return nil, err
	}
	message.Id = idMes

	return message, nil
}

func (r *repository) GetUsernameByUUID(ctx context.Context, uuid string) (string, error) {
	query := "SELECT username FROM users WHERE uuid = $1;"

	var username string
	err := r.DB.QueryRowContext(ctx, query, uuid).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}
