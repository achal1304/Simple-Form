package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gopheramit/Simple-Form/internal/validator"
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
	INSERT INTO users(email, apiKey)
	VALUES ($1, $2)
	RETURNING id, created_at, version`
	fmt.Println("in db")
	fmt.Println(user.Email)
	fmt.Println(user.ApiKey)
	args := []interface{}{user.Email, user.ApiKey}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	err := m.DB.PingContext(ctx)
	if err != nil {
		//return nil, err
		fmt.Println("ping")
		fmt.Println(err)
	}
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.Version)
	defer cancel()
	fmt.Println(err)
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
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
