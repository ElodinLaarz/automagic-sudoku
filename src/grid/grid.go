package grid

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
)

const (
	DefaultCellClass     = "cell"
	HighlightedCellClass = "cell-highlighted"

	StandardGridSize = 3
	NumGrids         = 2
)

type Cell struct {
	OriginIndex int
	Class       string // HTML class.
	Id          string // HTML id.
	Value       uint64
}

type Row struct {
	Cells []Cell
}
type Grid struct {
	Rows      []Row
	GridIndex int
	GridSize  int
}

func cellIdString(cellIndex, gridIndex int) string {
	return fmt.Sprintf("cell-%d-grid-%d", cellIndex, gridIndex)
}

func (g *Grid) shuffleRows(seed int64) {
	rand.New(rand.NewSource(seed)).Shuffle(len(g.Rows), func(i, j int) { g.Rows[i], g.Rows[j] = g.Rows[j], g.Rows[i] })
}

func makeGrid(gridSize, gridIndex int) []Row {
	var rows []Row
	numRows := int(gridSize * gridSize)
	numCols := numRows
	for rowIndex := 0; rowIndex < numRows; rowIndex++ {
		row := Row{}
		for colIndex := 0; colIndex < numCols; colIndex++ {
			originIndex := numCols*rowIndex + colIndex
			row.Cells = append(row.Cells, Cell{
				OriginIndex: originIndex,
				Class:       DefaultCellClass,
				Id:          cellIdString(originIndex, gridIndex),
				Value:       0,
			})
		}
		rows = append(rows, row)
	}
	return rows
}

func MakeGrids(NumGrids, gridSize int) []Grid {
	var grids []Grid
	for gridIndex := 0; gridIndex < NumGrids; gridIndex++ {
		g := Grid{
			Rows:      makeGrid(gridSize, gridIndex),
			GridSize:  gridSize,
			GridIndex: gridIndex + 1,
		}
		if gridIndex != 0 {
			g.shuffleRows(int64(gridIndex))
		}
		grids = append(grids, g)
	}
	return grids
}

func absoluteToRelative(index, sideLength int) (int, int) {
	row := index / sideLength
	col := index % sideLength
	return row, col
}

func NeighborCells(index, size int) map[string]bool {
	nbs := map[string]bool{}
	sideLength := size * size
	rowN, colN := absoluteToRelative(index, sideLength)
	for curAbsIndex := 0; curAbsIndex < sideLength*sideLength; curAbsIndex++ {
		r, c := absoluteToRelative(curAbsIndex, sideLength)
		// neighbors have the same row, column, or equivalence class
		// mod size in row AND column
		sameRow := r == rowN
		sameCol := c == colN
		sameBox := r/size == rowN/size && c/size == colN/size
		if sameRow || sameCol || sameBox {
			for curGrid := 0; curGrid < NumGrids; curGrid++ {
				nbs[cellIdString(curAbsIndex, curGrid)] = true
			}
		}
	}
	return nbs
}

var (
	updateCount = 0
)

func RenderGrid(w http.ResponseWriter, gridData map[string][]Grid) error {
	fmt.Printf("update #%d\n", updateCount)
	updateCount++
	tmpl, err := template.ParseFiles("./src/templates/grid.html")
	if err != nil {
		return fmt.Errorf("error parsing template: %s", err)
	}
	if err := tmpl.Execute(w, gridData); err != nil {
		return fmt.Errorf("error rendering grids: %s", err)
	}
	return nil
}
