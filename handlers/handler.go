package handlers

import (
	"fmt"
	"net/http"
)

func Handler() {
	fmt.Println(" HTTP handler initialized")

	// Register the route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "Hello, you've reached the server on port 8000!\n")
	})

}
