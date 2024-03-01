package handler

import (
	"fmt"
	"net/http"

	"automagic-sudoku/src/grid"
)

func GridHandler(w http.ResponseWriter, r *http.Request) {
	gridData := map[string][]grid.Grid{
		"grid": grid.MakeGrids(grid.NumGrids, grid.StandardGridSize),
	}
	if err := grid.RenderGrid(w, gridData); err != nil {
		fmt.Fprintf(w, "Error initially rendering grids: %s", err)
	}
}
