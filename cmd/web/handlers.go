package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gopheramit/Simple-Form/internal/pkg"
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
	app.render(w, r, "home.page.tmpl", nil)

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
