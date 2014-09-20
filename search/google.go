// Package search : google performs searches against the google search engine.
package search

/*
	http://ajax.googleapis.com/ajax/services/search/web?v=1.0&rsz=8&q=news+iraq
	Returns results in JSON

	{
		"responseData":	{
			"results":[
	  			{
				   "GsearchResultClass": "GwebSearch",
				   "unescapedUrl": "http://en.wikipedia.org/wiki/Paris_Hilton",
				   "url": "http://en.wikipedia.org/wiki/Paris_Hilton",
				   "visibleUrl": "en.wikipedia.org",
				   "cacheUrl": "http://www.google.com/search?q\u003dcache:TwrPfhd22hYJ:en.wikipedia.org",
				   "title": "\u003cb\u003eParis Hilton\u003c/b\u003e - Wikipedia, the free encyclopedia",
				   "titleNoFormatting": "Paris Hilton - Wikipedia, the free encyclopedia",
				   "content": "In 2006, she released her debut album \u003cb\u003eParis\u003c/b\u003e..."
	  			}
			]
		}
	}
*/

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Google provides support for Google searches.
type Google struct{}

type responseData struct {
	Results []map[string]interface{} `json:"results"`
}

type gResponse struct {
	ResponseData responseData `json:"responseData"`
}

// NewGoogle returns a Google Searcher value.
func NewGoogle() Searcher {
	return Google{}
}

// Search implements the Searcher interface. It performs a search
// against Google.
func (g Google) Search(searchTerm string) *Result {
	log.Printf("Google Search : Started : searchTerm[%s]\n", searchTerm)

	// Build the search url for the call.
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	uri := "http://ajax.googleapis.com/ajax/services/search/web?v=1.0&rsz=8&q=" + searchTerm
	log.Printf("Google Search : URL : %s\n", uri)

	// Create a result value.
	// I need an empty slice over a nil slice so I can
	// return an empty JSON document if necessary.
	result := Result{
		Engine:  "Google",
		Results: []map[string]interface{}{},
	}

	// Issue the search against Google.
	resp, err := http.Get(uri)
	if err != nil {
		log.Printf("Google Search : Get : ERROR : %s\n", err)
		result.Error = err
		return &result
	}

	// Schedule the close of the response body.
	defer resp.Body.Close()

	// Decode the results into the slice of maps.
	var gr gResponse
	err = json.NewDecoder(resp.Body).Decode(&gr)
	if err != nil {
		log.Printf("Google Search : Decode : ERROR : %s\n", err)
		result.Error = err
		return &result
	}

	result.Results = gr.ResponseData.Results

	return &result
}
