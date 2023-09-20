package auth

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
	NewAccessToken(ctx context.Context, idUser string, access string) (*session, error)
	GetSessionByID(ctx context.Context, idUser string) (*session, error)
	NewAccessAndRefreshToken(ctx context.Context, idUser string, access string, refresh string) (*session, error)
}

func (r *repository) NewAccessAndRefreshToken(ctx context.Context, idUser string, access string,
	refresh string) (*session, error) {
	query := "UPDATE jwt SET accesstoken = $1, refreshtoken = $2  WHERE uuid = $3 ;"

	res, err := r.DB.ExecContext(ctx, query, access, refresh, idUser)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		query = "INSERT INTO jwt (uuid, accesstoken, refreshtoken) VALUES ($1, $2, $3);"

		_, err = r.DB.Exec(query, idUser, access, refresh)
		if err != nil {
			r.logger.Error(err.Error())
			return nil, err
		}

	}
	u := &session{
		accessToken:  access,
		refreshToken: refresh,
		uuid:         idUser,
	}

	return u, nil

}

func (r *repository) NewAccessToken(ctx context.Context, idUser, access string) (*session, error) {
	query := "UPDATE jwt SET accesstoken = $1  WHERE uuid = $2 ;"

	_, err := r.DB.ExecContext(ctx, query, access, idUser)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, err
	}

	u := &session{
		accessToken:  access,
		refreshToken: "",
		uuid:         idUser,
	}

	return u, nil
}

func (r *repository) GetSessionByID(ctx context.Context, idUser string) (*session, error) {
	query := "SELECT * FROM jwt WHERE uuid = $1;"

	sess := &session{}

	err := r.DB.QueryRowContext(ctx, query, idUser).Scan(&sess.uuid, &sess.accessToken, &sess.refreshToken)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
