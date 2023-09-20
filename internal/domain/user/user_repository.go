package user

import (
	"WSChats/internal/adapters/api/DTO"
	"WSChats/pkg/logger"
	"context"
	"database/sql"
	"errors"
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
	CreateUser(ctx context.Context, user *User) (*DTO.CreateUserRes, error)
	GetUserByEmail(ctx context.Context, req *DTO.GetUserByEmailReq) (*User, error)
	//UpdateUser(ctx context.Context) (*User, error)
	//DeleteUser(ctx context.Context) error
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*DTO.CreateUserRes, error) {
	query := "SELECT * FROM users WHERE email = $1;"

	row := r.DB.QueryRowContext(ctx, query, user.Email)

	if errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, errors.New("user has been created")
	}

	var resUser = &DTO.CreateUserRes{}

	query = "INSERT INTO users (uuid, username, email, password) VALUES ($1, $2, $3, $4);"

	_, err := r.DB.Exec(query, user.UUID, user.Username, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	resUser.Username = user.Username
	resUser.UUID = user.UUID
	resUser.Email = user.Email

	return resUser, nil
}
func (r *repository) GetUserByEmail(ctx context.Context, req *DTO.GetUserByEmailReq) (*User, error) {
	query := "SELECT * FROM users WHERE email = $1;"

	var u User

	err := r.DB.QueryRowContext(ctx, query, req.Email).Scan(&u.UUID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
func (r *repository) UpdateUser(ctx context.Context) {

}
func (r *repository) DeleteUser(ctx context.Context) {

}
