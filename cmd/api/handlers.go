package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gopheramit/Simple-Form/internal/data"
	"github.com/gopheramit/Simple-Form/internal/validator"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	key := app.genUlid()
	keystr := key.String()
	user := &data.User{

		Email:  input.Email,
		ApiKey: keystr,
	}
	v := validator.New()

	if data.ValidateEmail(v, user.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Println(user.Email)
	fmt.Println(user.ApiKey)
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {

		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"apiKey": user.ApiKey}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
