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

func (m FormModel) KeyVerification(key string) (string, error) {
	query := ` select email FROM users WHERE apiKey= $1`
	var email string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, key).Scan(&email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", err

		}
	}
	return email, nil

}
