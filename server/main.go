package main

import (
	"html/template"
	"net/http"
)

type Todo struct {
	Title string
	Done  bool
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("../templates/index.html"))

		todos := []Todo{
			{Title: "Learn Go", Done: true},
			{Title: "Build Todo App", Done: false},
			{Title: "Deploy", Done: false},
		}

		data := map[string]interface{}{
			"Title": "Todo list",
			"Todos": todos,
		}
		tmpl.Execute(w, data)
	})

	// React用Cors
	// handler := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"http://localhost:3000"},
	// 	AllowedMethods: []string{http.MethodGet, http.MethodPost, http.Methodput, http.MethodDelete},
	// 	AllowedHeaders: []string{"*"},
	// }).Handler(mux)

	http.ListenAndServe(":8080", mux)
}
