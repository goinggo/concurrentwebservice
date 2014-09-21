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

	// Capture all the form values.
	fv := formValues(r)

	// If this is a post, perform a search.
	if r.Method == "POST" {
		if fv["searchterm"] != "" {
			results = search.Submit(fv)
		}
	}

	// Render the index page.
	view, err := renderIndex(fv, results)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	// Write the final markup as the response.
	fmt.Fprint(w, string(view))
}

// formValues extracts the form data.
func formValues(r *http.Request) map[string]interface{} {
	fv := make(map[string]interface{})
	fv["searchterm"] = r.FormValue("searchterm")

	if r.FormValue("google") == "on" {
		fv["google"] = "checked"
	} else {
		fv["google"] = ""
	}

	if r.FormValue("bing") == "on" {
		fv["bing"] = "checked"
	} else {
		fv["bing"] = ""
	}

	if r.FormValue("blekko") == "on" {
		fv["blekko"] = "checked"
	} else {
		fv["blekko"] = ""
	}

	if r.FormValue("first") == "on" {
		fv["first"] = "checked"
	} else {
		fv["first"] = ""
	}

	return fv
}

// renderIndex generates the HTML response for this route.
func renderIndex(fv map[string]interface{}, results []search.Result) ([]byte, error) {
	// Generate the HTML for the results content.
	if results != nil {
		html, err := renderResult(results)
		if err != nil {
			fv["Results"] = err.Error()
		}

		fv["Results"] = template.HTML(string(html))
	}

	// Generate the HTML for the index content.
	html := bytes.NewBufferString("")
	if err := views["index"].Execute(html, fv); err != nil {
		log.Printf("Index Service : Index : ERROR : %s\n", err)
		return nil, err
	}

	// Bind the layout markup for the final document.
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
