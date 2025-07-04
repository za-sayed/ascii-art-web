package main

import (
	"ascii-art-web/asciiart"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// Preloaded templates map
var templates = map[string]*template.Template{}

type Data struct {
	Text   string
	Result string
}

func init() {
	loadTemplates("home.html", "error.html")
}

func main() {
	// Route handlers
	http.HandleFunc("/", routeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)

	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))

	log.Println("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}

// loadTemplates parses the provided template files and stores them in the templates map
func loadTemplates(files ...string) {
	for _, file := range files {
		tmpl, err := template.ParseFiles("templates/" + file)
		if err != nil {
			log.Fatalf("Error loading template %s: %v", file, err)
		}
		templates[file] = tmpl
	}
}

// routeHandler handles different URL paths
func routeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		if r.Method == http.MethodGet {
			renderTemplate(w, "home.html", nil)
		} else {
			renderError(w, http.StatusMethodNotAllowed)
		}
	default:
		renderError(w, http.StatusNotFound)
	}
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderError(w, http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	style := r.FormValue("style")

	if strings.TrimSpace(text) == "" || !asciiart.Validation(text) {
		renderError(w, http.StatusBadRequest)
		return
	}

	if style != "standard" && style != "thinkertoy" && style != "shadow" {
		renderError(w, http.StatusBadRequest)
		return
	}

	// Load banner and generate ASCII art
	Lresult, err := asciiart.Load("banners\\" + style + ".txt")
	if err != nil {
		renderError(w, http.StatusInternalServerError)
		return
	}
	result := asciiart.GenerateAsciiArt(text, Lresult)

	data := Data{
		Text:   text,
		Result: result,
	}

	renderTemplate(w, "home.html", data)
}

// renderTemplate renders the specified template with provided data
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, ok := templates[tmpl]
	if !ok {
		renderError(w, http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		renderError(w, http.StatusInternalServerError)
	}
}

// renderError renders an error page based on the status code
func renderError(w http.ResponseWriter, status int) {
	data := struct {
		Status  int
		Message string
	}{
		Status:  status,
		Message: http.StatusText(status),
	}
	w.WriteHeader(status)
	renderTemplate(w, "error.html", data)
}
