package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type FormModel struct {
	DB *sql.DB
}
type Form struct {
	AccessKey string `json:"access-key"`
	Name      string `json:"name"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Message   string `json:"message"`
}

func (m FormModel) KeyVerification(key string) (string, int64, error) {
	query := ` select email,apiCount FROM users WHERE apiKey= $1`
	var email string
	var count int64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, key).Scan(&email, &count)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", 0, ErrRecordNotFound
		default:
			return "", 0, err

		}
	}
	return email, count, nil

}

func (m FormModel) UpdateCount(key string, cnt int64) error {
	query := `update users set apiCount=$1 WHERE apiKey= $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, cnt, key)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err

		}
	}
	return nil

}
