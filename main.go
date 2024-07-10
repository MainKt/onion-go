package main

import (
	"html/template"
	"log"
	"time"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"math/rand"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	phones := []string{ "86184 97080", "91486 55749", "78210 99805" };

	promptMessages := []string{ "See carefully, LISTEN carefully" };

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))

		if err := tmpl.Execute(w, struct {
			InvalidPassword bool
			Phone string
			Prompt string
		}{
			InvalidPassword: r.Header.Get("Referer") != "",
			Phone: phones[rand.Intn(len(phones))],
			Prompt: promptMessages[0],
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

		if err := tmpl.Execute(w, struct { Phone string }{ Phone: phones[rand.Intn(len(phones))] }); err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})

	log.Println("Listening on :4104...")
	err := http.ListenAndServe(":4104", nil)
	if err != nil {
		log.Fatal(err)
	}
}
