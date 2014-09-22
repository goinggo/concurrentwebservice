// Package search manages the searching of results against Google, Yahoo and Bing.
package search

import (
	"html/template"
	"log"
	"sync"
)

// Options provides the search options for performing searches.
type Options struct {
	SearchTerm string
	Google     bool
	Bing       bool
	Blekko     bool
	First      bool
}

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
func Submit(options *Options) []Result {
	log.Printf("search : Submit : Started : %#v\n", options)

	var final []Result
	searchers := make(map[string]Searcher)
	searchResults := make(chan []Result)

	// Create a Google Searcher if checked.
	if options.Google {
		log.Println("search : Submit : Info : Adding Google")
		searchers["google"] = NewGoogle()
	}

	// Create a Bing Searcher if checked.
	if options.Bing {
		log.Println("search : Submit : Info : Adding Bing")
		searchers["bing"] = NewBing()
	}

	// Create a Bing Searcher if checked.
	if options.Blekko {
		log.Println("search : Submit : Info : Adding Blekko")
		searchers["blekko"] = NewBlekko()
	}

	// Perform the searches concurrently. Using a map because
	// it returns the searchers in a random order every time.
	for _, searcher := range searchers {
		go searcher.Search(options.SearchTerm, searchResults)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		var sent bool

		for search := 0; search < len(searchers); search++ {
			// Wait to recieve results.
			log.Println("search : Submit : Info : Waiting For Results...")
			sr := <-searchResults

			if sent {
				continue
			}

			// Save the results to the results slice.
			log.Printf("search : Submit : Info : Results Returned : Results[%d]\n", len(sr))
			final = append(final, sr...)

			// If we just want the first result, don't wait any longer and give
			// the user the results we have.
			if options.First {
				sent = true
				wg.Done()
			}
		}

		if !sent {
			wg.Done()
		}

		log.Println("search : Submit : Info : All Results Are In")
	}()

	wg.Wait()

	log.Printf("search : Submit : Completed : Found [%d] Results\n", len(final))
	return final
}
