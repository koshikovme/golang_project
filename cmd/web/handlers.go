package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Define a home handler func which writes a byte slice containing
// "Hello from Snippetbox" as resp body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the curr req path exactly matches "/". If it doesn't, use
	// the http.NotFound() func to send 404 resp.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		// return from func to avoid proceeding to home page response
		return
	}

	files := []string{
		"./ui/html/home.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}
	// Use the template.ParseFiles() func to read the files and store the templates
	// in a template set. Noticed that we can pass teh slice of file paths as a variadic
	// param?
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// We then use the Execute() method on the template set to write the template
	// content as the response body. The last parameter to Execute() represents any
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract id value from query string and try to convert it to an integer
	// using strconv.Atoi() func. If it can't be converted, or the value is less than 1,
	// we return a 404 page Not Found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	if _, err := fmt.Fprintf(w, "Display a specific snippet with ID %d", id); err != nil {
		app.infoLog.Println("show snippet request:", err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// only allow Post
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", 405)
		return
	}

	if _, err := w.Write([]byte("Create a new snippet")); err != nil {
		app.infoLog.Println("create new snippet request:", err)
	}
}
