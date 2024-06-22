package http

import "net/http"

func RunServer() {
	// Start the server
	http.ListenAndServe(":8080", nil)
}
