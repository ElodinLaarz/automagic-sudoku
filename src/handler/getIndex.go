package handler

import (
	"fmt"
	"net/http"
	"text/template"

	"automagic-sudoku/src/grid"
)

const (
	indexTmpl   = "src/templates/index.html"
	sidebarTmpl = "src/templates/sidebar.html"
	iconTmpl    = "src/templates/icon.html"
	mainTmpl    = "src/templates/main.html"
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
	tmpl, err := template.ParseFiles(indexTmpl, sidebarTmpl, iconTmpl, mainTmpl)
	if err != nil {
		fmt.Println("ParseFiles:", err)
		return
	}
	if err := tmpl.Execute(w, templateData); err != nil {
		fmt.Println("Execute:", err)
	}
}
