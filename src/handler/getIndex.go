package handler

import (
	"fmt"
	"net/http"
	"text/template"

	"automagic-sudoku/src/grid"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	grids := map[string][]grid.Grid{
		"grid": grid.MakeGrids(grid.NumGrids, grid.StandardGridSize),
	}
	tmpl, err := template.ParseFiles("src/templates/index.html", "src/templates/sidebar.html")
	if err != nil {
		fmt.Println("ParseFiles:", err)
		return
	}
	if err := tmpl.Execute(w, grids); err != nil {
		fmt.Println("Execute:", err)
	}
}
