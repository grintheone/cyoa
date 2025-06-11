package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "path: %s", r.URL.Path)
}

func main() {
	port := ":8080"

	http.HandleFunc("/", handler)

	fmt.Printf("Starting a server on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
