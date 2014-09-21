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
func Submit(searchTerm string, google string, bing string, first string) []Result {
	log.Printf("Google Search : Started : searchTerm[%s] google[%s] bing[%s] first[%s]\n", searchTerm, google, bing, first)

	var searchers []Searcher
	var final []Result
	searchResults := make(chan []Result)

	// Create a Google Searcher if checked.
	if google == "checked" {
		log.Println("Google Search : Info : Adding Google")
		searchers = append(searchers, NewGoogle())
	}

	// Create a Bing Searcher if checked.
	if bing == "checked" {
		log.Println("Google Search : Info : Adding Bing")
		searchers = append(searchers, NewBing())
	}

	// Perform the searches concurrently.
	for _, searcher := range searchers {
		go searcher.Search(searchTerm, searchResults)
	}

	// Wait for the results to come back.
	for search := 0; search < len(searchers); search++ {
		// Wait to recieve results and save them.
		final = append(final, <-searchResults...)

		// If we just want the first result, don't wait any longer.
		if first == "checked" {
			break
		}
	}

	log.Printf("Google Search : Completed : Found [%d] Results\n", len(final))

	return final
}
