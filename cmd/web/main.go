package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/grintheone/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "port to run the server on")
	fname := flag.String("name", "gophers.json", "supply the name of the file")
	flag.Parse()

	f, err := os.Open(*fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	story, err := cyoa.JsonStory(f)
	if err != nil {
		log.Fatal(err)
	}

	// custom tempalte to show how functional options
	// allow us to customize our application
	tmpl := template.Must(template.New("").Parse("Hello"))
	h := cyoa.NewHandler(story, cyoa.WithTemplate(tmpl), cyoa.WithCustomPathFn(customPathFn))

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./ui/static"))

	// default behaviour
	mux.Handle("/", cyoa.NewHandler(story))
	// custom behaviour
	mux.Handle("/story/", h)
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("server is listening on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func customPathFn(r *http.Request) string {
	path := r.URL.Path
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}

	return path[len("/story/"):]
}
