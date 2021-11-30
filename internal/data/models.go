package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users UserModel
	Forms FormModel
}

func NewModels(db *sql.DB) Models {
	return Models{

		Users: UserModel{DB: db},
		Forms: FormModel{DB: db},
	}
}
