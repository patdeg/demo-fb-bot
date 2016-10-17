// Demo of Facebook Bot with Google App Engine
package main

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
)

// HTML Template for the home page
var homeTemplate = template.Must(template.New("index.html").Delims("[[", "]]").ParseFiles("index.html"))

// Main init function to assign paths to handlers
func init() {

	// Home page (& catch-all)
	http.HandleFunc("/", HomeHandler)

	// Facebook Callback
	http.HandleFunc("/callback", FacebookCallbackHandler)

}

// Home page handler
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	log.Debugf(c, ">>> Home Handler")

	// Execute the home page template
	if err := homeTemplate.Execute(w, template.FuncMap{
		"Version": appengine.VersionID(c),
	}); err != nil {
		log.Errorf(c, "Error with homeTemplate: %v", err)
		http.Error(w, "Internal Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
