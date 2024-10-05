package handler

import "net/http"

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	// Return a JSON message with the API name
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Welcome to the Quiz API"}`))
}
