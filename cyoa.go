package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"path"
	"strings"
)

const templateName = "base.html"

var tpl *template.Template

func init() {
	path := path.Join("ui", "html", templateName)
	tpl = template.Must(template.New(templateName).Funcs(functions).ParseFiles(path))
}

func isEmpty(slice []Option) bool {
	return len(slice) == 0
}

var functions = template.FuncMap{
	"isEmpty": isEmpty,
}

type Story map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithCustomPathFn(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path[1:])
	if path == "" || path == "/" {
		path = "intro"
	}
	return path
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found", http.StatusNotFound)
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)

	story := make(Story)

	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}
