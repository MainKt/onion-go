package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))

		if err := tmpl.Execute(w, struct {
			InvalidPassword bool
		}{
			InvalidPassword: r.Header.Get("Referer") != "",
		}); err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("POST /research", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		password := r.Form.Get("password")
		if password != os.Getenv("PASSWD") {
			http.Redirect(w, r, "/#password", http.StatusSeeOther)
			return
		}

		tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "research.html")))

		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})

	log.Println("Listening on :4104...")
	err := http.ListenAndServe(":4104", nil)
	if err != nil {
		log.Fatal(err)
	}
}
