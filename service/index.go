// Package service : index maintains the support for the home page.
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/goinggo/concurrentwebservice/search"
)

// index handles the home page route processing.
func index(w http.ResponseWriter, r *http.Request) {
	var result *search.Result
	if r.Method == "POST" {
		result = search.Submit("news iraq")
	}

	view, err := renderIndex(result)
	if err != nil {
		fmt.Fprint(w, err)
	}

	fmt.Fprint(w, string(view))
}

// renderIndex generates the HTML response for this route.
func renderIndex(result *search.Result) ([]byte, error) {
	vars := make(map[string]interface{})
	if result != nil {
		data, err := json.MarshalIndent(result.Results, "", "    ")
		if err != nil {
			vars["Results"] = err.Error()
		}
		vars["Results"] = string(data)
	}

	// Generate the HTML for the index content.
	index := bytes.NewBufferString("")
	if err := views["index"].Execute(index, vars); err != nil {
		log.Println(err)
		return nil, err
	}

	return renderLayout(index)
}
