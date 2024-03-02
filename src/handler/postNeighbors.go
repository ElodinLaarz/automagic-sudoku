package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"automagic-sudoku/src/grid"
)

func NeighborHandler(w http.ResponseWriter, r *http.Request) {
	ind := r.URL.Query().Get("originIndex")
	ind = strings.Trim(ind, "[]")
	originInd := strings.Split(ind, " ")
	fmt.Printf("originInd: %v\n", originInd)
	originIndicies := make([]int, len(originInd))
	for i, v := range originInd {
		if v == "" {
			continue
		}
		ix, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(w, "Error converting %s to int: %s", v, err)
			return
		}
		originIndicies[i] = ix
	}
	relativeIndex, err := strconv.Atoi(r.URL.Query().Get("relativeIndex"))
	if err != nil {
		fmt.Fprintf(w, "Error converting %s to int: %s", r.URL.Query().Get("relativeIndex"), err)
		return
	}
	gridName := r.URL.Query().Get("gridName")
	neighbors := grid.NeighborCells(relativeIndex, grid.SingleBoxSize, gridName)

	for _, g := range grid.Grids {
		for rowIndex, row := range g.Rows {
			for colIndex, cell := range row.Cells {
				if _, ok := neighbors[cell.Id]; ok {
					g.Rows[rowIndex].Cells[colIndex].Class = grid.HighlightedCellClass
				} else {
					g.Rows[rowIndex].Cells[colIndex].Class = grid.DefaultCellClass
				}
			}
		}
	}
	gridData := make(map[string]map[string]grid.Grid)
	gridData["grid"] = grid.Grids
	if err := grid.RenderGrid(w, gridData); err != nil {
		fmt.Fprintf(w, "Error highlighting id %d: %s", originIndicies, err)
	}
}
