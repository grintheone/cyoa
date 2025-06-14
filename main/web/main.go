package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	var (
		port  = ":8080"
		paths = make(PathHandler)
	)

	f, err := os.Open("gophers.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&paths); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", paths)

	fs := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Starting a server on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
