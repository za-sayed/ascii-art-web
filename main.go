package main

import (
	"ascii-art-web/asciiart"
	"html/template"
	"net/http"
	"log"
)

// Preloaded templates map
var templates = map[string]*template.Template{}

type Data struct {
	Text   string
	Result string
}

func init() {
	loadTemplates("home.html", "400.html", "404.html", "500.html")
}


func main() {
	// Route handlers
	http.HandleFunc("/", routeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)

	// Serve static files from the static directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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

	if text == "" || !asciiart.Validation(text) {
		renderError(w, http.StatusBadRequest)
		return
	}

	if style == "" {
		style = "standard"
	}

	// Load banner and generate ASCII art
	load, err := asciiart.Load(style + ".txt")
	if err != nil {
		renderError(w, http.StatusInternalServerError)
		return
	}
	result := asciiart.GenerateAsciiArt(text, load)

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
	var errorTemplate string
	switch status {
	case http.StatusBadRequest:
		errorTemplate = "400.html"
	case http.StatusNotFound:
		errorTemplate = "404.html"
	case http.StatusInternalServerError:
		errorTemplate = "500.html"
	default:
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.WriteHeader(status)
	renderTemplate(w, errorTemplate, nil)
}
