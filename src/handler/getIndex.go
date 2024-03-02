package handler

import (
	"fmt"
	"net/http"
	"text/template"

	"automagic-sudoku/src/grid"
)

const (
	index   = "src/templates/index.html"
	sidebar = "src/templates/sidebar.html"
	icon    = "src/templates/icon.html"
)

type IconTemplate struct {
	Icon string
}

var exampleIcons = []IconTemplate{
	{Icon: "icon1"},
	{Icon: "icon2"},
	{Icon: "icon3"},
	{Icon: "icon4"},
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	templateData := map[string]any{
		"grid": grid.MakeGrids(grid.NumGrids, grid.StandardGridSize),
		"icon": exampleIcons,
	}
	tmpl, err := template.ParseFiles(index, sidebar, icon)
	if err != nil {
		fmt.Println("ParseFiles:", err)
		return
	}
	if err := tmpl.Execute(w, templateData); err != nil {
		fmt.Println("Execute:", err)
	}
}
