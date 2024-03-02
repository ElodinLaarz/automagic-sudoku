package handler

import (
	"fmt"
	"net/http"
	"text/template"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(mainTmpl)
	if err != nil {
		fmt.Println("ParseFiles:", err)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println("Execute:", err)
	}
}
