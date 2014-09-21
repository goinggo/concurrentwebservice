// Package service : temnplates provides support for using HTML
// based templates for responses.
package service

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
)

// views contains a map of templates for rendering views.
var views = make(map[string]*template.Template)

// init loads the existing templates for use by routing code.
func init() {
	loadTemplate("layout", "views/basic-layout.html")
	loadTemplate("index", "views/index.html")
	loadTemplate("results", "views/results.html")
}

// loadTemplate reads the specified template file for use.
func loadTemplate(name string, path string) {
	// Read the html template file.
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	// Create a template value for this code.
	tmpl, err := template.New(name).Parse(string(data))
	if err != nil {
		log.Fatalln(err)
	}

	// Have we processed this file already?
	if _, exists := views[name]; exists {
		log.Fatalf("Template %s already in use.", name)
	}

	// Store the template for use.
	views[name] = tmpl
}

// renderLayout generates the HTML response with the layout
func renderLayout(content []byte) ([]byte, error) {
	// Place the layout content into a map for processing.
	vars := make(map[string]interface{})
	vars["LayoutContent"] = template.HTML(string(content))

	// Generate the final markup by embedding the index content
	// into the layout markup.
	final := new(bytes.Buffer)
	if err := views["layout"].Execute(final, vars); err != nil {
		log.Println(err)
		return nil, err
	}

	// Return the final markup for the reponse.
	return final.Bytes(), nil
}
