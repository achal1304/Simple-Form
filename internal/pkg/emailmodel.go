package pkg

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type EmailInput struct {
	Email string `json:"email"`
}
