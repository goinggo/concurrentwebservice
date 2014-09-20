// Package service : index maintains the support for the home page.
package service

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// index handles the home page route processing.
func index(w http.ResponseWriter, r *http.Request) {
	result := ""
	if r.Method == "POST" {
		result = "Comming Soon"
	}

	view, err := renderIndex(result)
	if err != nil {
		fmt.Fprint(w, err)
	}

	fmt.Fprint(w, string(view))
}

// renderIndex generates the HTML response for this route.
func renderIndex(result string) ([]byte, error) {
	vars := make(map[string]interface{})
	vars["Results"] = result

	// Generate the HTML for the index content.
	index := bytes.NewBufferString("")
	if err := views["index"].Execute(index, vars); err != nil {
		log.Println(err)
		return nil, err
	}

	return renderLayout(index)
}
