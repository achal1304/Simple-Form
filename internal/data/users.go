package data

import (
	"database/sql"
	"time"
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
