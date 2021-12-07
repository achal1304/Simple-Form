package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gopheramit/Simple-Form/internal/pkg"
	"github.com/gopheramit/Simple-Form/internal/pkg/forms"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }
	//s, err := app.snippets.Latest()
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	app.render(w, r, "home.page.tmpl", &templateData{
		Form: forms.New(nil)})

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	response, err := http.Get("http://localhost:4000/v1/healthcheck")
	var input *pkg.EmailInput

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))
	if err := json.Unmarshal(responseData, &input); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	app.render(w, r, "show.page.tmpl", &templateData{Snippet: input})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil)})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)

	// 	return
	// }
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email")
	form.MaxLength("email", 100)
	// form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	httpposturl := "http://localhost:4000/v1/users"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var emailForm = form.Get("email")

	type myStruct struct {
		Email string `json:"email"`
	}

	myData := myStruct{
		Email: emailForm,
	}

	jsonData, err := json.Marshal(myData)

	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}

	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	// id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// app.session.Put(r, "flash", "Snippet successfully created!")

	// http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
