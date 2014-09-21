// Package search manages the searching of results against Google, Yahoo and Bing.
package search

import (
	"html/template"
	"log"
)

// Result represents a search result that was found.
type Result struct {
	Engine  string
	Title   string
	Link    string
	Content string
}

// TitleHTML fixes encoding issues.
func (r *Result) TitleHTML() template.HTML {
	return template.HTML(r.Title)
}

// ContentHTML fixes encoding issues.
func (r *Result) ContentHTML() template.HTML {
	return template.HTML(r.Content)
}

// Searcher declares an interface used to leverage different
// search engines to find results.
type Searcher interface {
	Search(searchTerm string, searchResults chan<- []Result)
}

// Submit uses goroutines and channels to perform a search against the three
// leading search engines concurrently.
func Submit(fv map[string]interface{}) []Result {
	log.Printf("Google Search : Started : %#v\n", fv)

	var searchers []Searcher
	var final []Result
	searchResults := make(chan []Result)

	// Create a Google Searcher if checked.
	if fv["google"].(string) == "checked" {
		log.Println("Submit : Info : Adding Google")
		searchers = append(searchers, NewGoogle())
	}

	// Create a Bing Searcher if checked.
	if fv["bing"].(string) == "checked" {
		log.Println("Submit : Info : Adding Bing")
		searchers = append(searchers, NewBing())
	}

	// Create a Bing Searcher if checked.
	if fv["blekko"].(string) == "checked" {
		log.Println("Submit : Info : Adding Blekko")
		searchers = append(searchers, NewBlekko())
	}

	// Perform the searches concurrently.
	for _, searcher := range searchers {
		go searcher.Search(fv["searchterm"].(string), searchResults)
	}

	// Wait for the results to come back.
	for search := 0; search < len(searchers); search++ {
		// Wait to recieve results and save them.
		log.Println("Submit : Info : Waiting For Results...")
		results := <-searchResults

		// Save the results to the final slice.
		log.Printf("Submit : Info : Results Returned : Results[%d]\n", len(results))
		final = append(final, results...)

		// If we just want the first result, don't wait any longer.
		if fv["first"].(string) == "checked" {
			break
		}
	}

	log.Printf("Submit : Completed : Found [%d] Results\n", len(final))

	return final
}
