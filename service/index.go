// Package service : index maintains the support for the home page.
package service

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/goinggo/concurrentwebservice/search"
)

// index handles the home page route processing.
func index(w http.ResponseWriter, r *http.Request) {
	var results []search.Result

	vars := variables(r)

	if r.Method == "POST" {
		if vars["searchTerm"] != "" {
			results = search.Submit(vars["searchTerm"].(string), vars["google"].(string), vars["bing"].(string), vars["first"].(string))
		}
	}

	view, err := renderIndex(vars, results)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, string(view))
}

// variables extracts the form data.
func variables(r *http.Request) map[string]interface{} {
	vars := make(map[string]interface{})
	vars["searchTerm"] = r.FormValue("searchterm")

	if r.FormValue("google") == "on" {
		vars["google"] = "checked"
	} else {
		vars["google"] = ""
	}

	if r.FormValue("bing") == "on" {
		vars["bing"] = "checked"
	} else {
		vars["bing"] = ""
	}

	if r.FormValue("first") == "on" {
		vars["first"] = "checked"
	} else {
		vars["first"] = ""
	}

	return vars
}

// renderIndex generates the HTML response for this route.
func renderIndex(vars map[string]interface{}, results []search.Result) ([]byte, error) {
	// Generate the HTML for the results content.
	if results != nil {
		html, err := renderResult(results)
		if err != nil {
			vars["Results"] = err.Error()
		}

		vars["Results"] = template.HTML(string(html))
	}

	// Generate the HTML for the index content.
	html := bytes.NewBufferString("")
	if err := views["index"].Execute(html, vars); err != nil {
		log.Printf("Index Service : Index : ERROR : %s\n", err)
		return nil, err
	}

	return renderLayout(html.Bytes())
}

// renderResult produces the HTML for the results.
func renderResult(items []search.Result) ([]byte, error) {
	vars := make(map[string]interface{})
	vars["Items"] = items

	// Generate the HTML for the index content.
	html := bytes.NewBufferString("")
	if err := views["results"].Execute(html, vars); err != nil {
		log.Printf("Index Service : Results : ERROR : %s\n", err)
		return nil, err
	}

	return html.Bytes(), nil
}
