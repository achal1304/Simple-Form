package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/oklog/ulid"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0,
			errors.New("invalid id parameter")
	}
	return id, nil
}

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	maxBytes := 1_048_576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}

	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil

}

func (app *application) genUlid() ulid.ULID {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id
}

func (app *application) SetupMail(email string, apiKey string) error {
	from := app.config.smtp.username

	password := app.config.smtp.password

	host := app.config.smtp.host
	port := app.config.smtp.port
	portStr := strconv.Itoa(port)
	toList := []string{email}

	auth := smtp.PlainAuth("", from, password, host)

	t, _ := template.ParseFiles("internal/html/welcome.tmpl")
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Access key from SIMPLE FORMS \n%s\n\n", mimeHeaders)))
	t.Execute(&body, struct {
		Email     string
		AccessKey string
	}{
		Email:     email,
		AccessKey: apiKey,
	})

	err := smtp.SendMail(host+":"+portStr, auth, from, toList, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Successfully sent mail to all user in toList")
	return nil
}

func (app *application) SendFormOnMail(email string, data interface{}) error {
	from := app.config.smtp.username

	password := app.config.smtp.password

	host := app.config.smtp.host
	port := app.config.smtp.port
	portStr := strconv.Itoa(port)
	toList := []string{email}
	// msg := "Hello geeks!!!"
	// body := []byte(msg)
	auth := smtp.PlainAuth("", from, password, host)

	t, _ := template.ParseFiles("internal/html/form-data.tmpl")
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Query Form Details\n%s\n\n", mimeHeaders)))
	t.Execute(&body, struct {
		Email string
		Data  interface{}
	}{
		Email: email,
		Data:  data,
	})

	err := smtp.SendMail(host+":"+portStr, auth, from, toList, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Successfully sent mail to all user in toList")
	return nil
}
