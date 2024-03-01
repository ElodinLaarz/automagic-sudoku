package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"automagic-sudoku/src/grid"
)

func NeighborHandler(w http.ResponseWriter, r *http.Request) {
	originIndex, err := strconv.Atoi(r.URL.Query().Get("originIndex"))
	if err != nil {
		fmt.Fprintf(w, "Error parsing id: %s", err)
		return
	}
	neighbors := grid.NeighborCells(originIndex, grid.StandardGridSize)

	grids := grid.MakeGrids(grid.NumGrids, grid.StandardGridSize)
	for _, data := range grids {
		for rowIndex, row := range data.Rows {
			for colIndex, cell := range row.Cells {
				if _, ok := neighbors[cell.Id]; ok {
					data.Rows[rowIndex].Cells[colIndex].Class = grid.HighlightedCellClass
					data.Rows[rowIndex].Cells[colIndex].Value = 1
				}
			}
		}
	}
	gridData := make(map[string][]grid.Grid)
	gridData["grid"] = grids
	if err := grid.RenderGrid(w, gridData); err != nil {
		fmt.Fprintf(w, "Error highlighting id %d: %s", originIndex, err)
	}
}
