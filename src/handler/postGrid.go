package handler

import (
	"fmt"
	"net/http"

	"automagic-sudoku/src/grid"
)

func GridHandler(w http.ResponseWriter, r *http.Request) {
	reset := r.URL.Query().Get("reset")
	fmt.Printf("reset: %s\n", reset)
	if reset == "true" {
		for _, g := range grid.Grids {
			g.DefaultClass()
		}
	}
	gridData := map[string]map[string]grid.Grid{
		"grid": grid.Grids,
	}
	if err := grid.RenderGrid(w, gridData); err != nil {
		fmt.Fprintf(w, "Error initially rendering grids: %s", err)
	}
}
