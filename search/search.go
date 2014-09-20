// Package search manages the searching of results against Google, Yahoo and Bing.
package search

import (
	"time"
)

// Result contains the search result
type Result struct {
	Engine   string
	Results  []map[string]interface{}
	Order    int
	Duration time.Time
	Error    error
}

// Searcher declares an interface used to leverage different
// search engines to find results.
type Searcher interface {
	Search(searchTerm string) *Result
}

// Submit uses goroutines and channels to perform a search against the three
// leading search engines concurrently.
func Submit(searchTerm string) *Result {
	gResult := NewGoogle().Search(searchTerm)

	return gResult
}
