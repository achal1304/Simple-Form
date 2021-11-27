package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserModel struct {
	DB *sql.DB
}
type User struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	ApiKey    string    `json:"apiKey"`
	ApiCount  int64     `json:"apiCount"`
	CreatedAt time.Time `json:"created_at"`
	Version   int64     `json:"version"`
}

func (m UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (id, email, apiKey, apiCount,created_at,version)
	VALUES ($1, $2, $3, $4,$5,$6)
	RETURNING id, created_at, version`
	args := []interface{}{user.Id, user.Email, user.ApiKey, user.ApiCount, user.CreatedAt, user.Version}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}
