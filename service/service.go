// Package service maintains the logic for the web service.
package service

import (
	"net/http"
)

// init binds the routes and handlers for the web service.
func init() {
	http.HandleFunc("/", index)
}

// Start binds the service to a port and starts listening
// for requests.
func Start() {
	http.ListenAndServe("localhost:9999", nil)
}
