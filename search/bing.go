// Package search : bing performs searches against the bing search engine.
package search

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"
)

// Bing provides support for Bing searches.
type Bing struct{}

type (
	// Item defines the fields associated with the item tag in the buoy RSS document.
	Item struct {
		XMLName     xml.Name `xml:"item"`
		PubDate     string   `xml:"pubDate"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
		Link        string   `xml:"link"`
	}

	// Image defines the fields associated with the image tag in the buoy RSS document.
	Image struct {
		XMLName xml.Name `xml:"image"`
		URL     string   `xml:"url"`
		Title   string   `xml:"title"`
		Link    string   `xml:"link"`
	}

	// Channel defines the fields associated with the channel tag in the buoy RSS document.
	Channel struct {
		XMLName xml.Name `xml:"channel"`
		Image   Image    `xml:"image"`
		Items   []Item   `xml:"item"`
	}

	// Document defines the fields associated with the buoy RSS document.
	Document struct {
		XMLName xml.Name `xml:"rss"`
		Channel Channel  `xml:"channel"`
	}
)

// NewBing returns a Bing Searcher value.
func NewBing() Searcher {
	return Bing{}
}

// Search implements the Searcher interface. It performs a search
// against Bing.
func (b Bing) Search(searchTerm string, searchResults chan<- []Result) {
	log.Printf("Bing Search : Started : searchTerm[%s]\n", searchTerm)

	// Need an empty slice so I can return an empty
	// JSON document if necessary.
	results := []Result{}

	// On return send the results we have.
	defer func() {
		searchResults <- results
	}()

	// Build a proper search url.
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	uri := "http://www.bing.com/search?q=" + searchTerm + "&format=rss"
	log.Printf("Bing Search : URL : %s\n", uri)

	// Issue the search against Google.
	resp, err := http.Get(uri)
	if err != nil {
		log.Printf("Bing Search : Get : ERROR : %s\n", err)
		return
	}

	// Schedule the close of the response body.
	defer resp.Body.Close()

	// Decode the results into the slice of maps.
	var d Document
	err = xml.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		log.Printf("Bing Search : Decode : ERROR : %s\n", err)
		return
	}

	// Capture the data we need for our results.
	for _, result := range d.Channel.Items {
		results = append(results, Result{
			Engine:  "Bing",
			Title:   result.Title,
			Link:    result.Link,
			Content: result.Description,
		})
	}

	log.Println("Bing Search : Completed")
}
