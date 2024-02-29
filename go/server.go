package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func main() {
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		_ = tmpl.Execute(w, nil)
	}
	h2 := func(w http.ResponseWriter, r *http.Request) {
		res := map[string]interface{}{
			"Name":  "Wyndham",
			"Phone": "8888888",
			"Email": "skyscraper@gmail.com",
		}
		tmpl := template.Must(template.ParseFiles("name_card.html"))
		_ = tmpl.Execute(w, res)
	}
	http.HandleFunc("/", h1)
	http.HandleFunc("/action", h2)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
