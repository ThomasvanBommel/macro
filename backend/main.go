package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to macro!")
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}