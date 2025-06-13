package main

import (
	"net/http"
)

type (
	Path        string
	PathHandler map[Path]Story
)

type Story struct {
	Title       string   `json:"title"`
	Description []string `json:"story"`
	Options     []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func (s Story) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := renderTemplate(w, s)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func (ph PathHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := Path(r.URL.Path[1:])
	if p == "" {
		p = "intro"
	}

	story, ok := ph[p]
	if !ok {
		http.Error(w, "Chapter not found.", http.StatusNotFound)
		return
	}

	story.ServeHTTP(w, r)
}
