package progress_tracker

import (
	"fmt"
	"net/http"
)

func Start() {
	// Initialize a server
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received from", r.RemoteAddr, "for", r.URL, "with method", r.Method)

		fmt.Fprintf(w, "Hello, world!")
	})
	http.ListenAndServe(":8080", nil)
}
