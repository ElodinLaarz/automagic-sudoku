package grid

import (
	"fmt"
	"html/template"
	"net/http"
)

const (
	DefaultCellClass     = "cell"
	HighlightedCellClass = "cell-highlighted"
	HighlightedMainClass = "cell-highlighted-main"

	DomainGridName = "Domain"

	SingleBoxSize = 2
	NumGrids      = 6
)

type Cell struct {
	OriginIndex   []int
	RelativeIndex int
	Class         string // HTML class.
	Id            string // HTML id.
	IsMain        bool
	// Value         []uint64 // To represent maps, there may be multiple values.
}

type Row struct {
	Cells []Cell
}
type Grid struct {
	Rows      []Row
	GridName  string
	BoxSize   int
	GridIndex int
	Tab       int
}

var (
	Grids = MakeGrids(NumGrids, SingleBoxSize)
)

func cellIdString(cellIndex, gridIndex int) string {
	return fmt.Sprintf("cell-%d-grid-%d", cellIndex, gridIndex)
}

func (g *Grid) DefaultClass() {
	for rowIndex := range g.Rows {
		for colIndex := range g.Rows[rowIndex].Cells {
			g.Rows[rowIndex].Cells[colIndex].Class = DefaultCellClass
			g.Rows[rowIndex].Cells[colIndex].IsMain = false
		}
	}
}

func (g *Grid) shuffle(m SudokuMap) error {
	if len(g.Rows) == 0 || len(g.Rows[0].Cells) == 0 {
		return fmt.Errorf("grid has no cells")
	}
	preimage := m.Preimage()

	newRows := []Row{}
	sideLength := g.BoxSize * g.BoxSize
	for rowIndex := range sideLength {
		newRowCells := []Cell{}
		for colIndex := range sideLength {
			relativeIndex := rowIndex*sideLength + colIndex
			newCell := Cell{
				OriginIndex:   preimage[relativeIndex],
				RelativeIndex: relativeIndex,
				Class:         DefaultCellClass,
				Id:            cellIdString(relativeIndex, g.GridIndex),
			}
			newRowCells = append(newRowCells, newCell)
		}
		newRows = append(newRows, Row{Cells: newRowCells})
	}
	g.Rows = newRows
	return nil
}

func makeGrid(singleBoxSize, gridIndex int) []Row {
	var rows []Row
	numRows := int(singleBoxSize * singleBoxSize)
	numCols := numRows
	for rowIndex := 0; rowIndex < numRows; rowIndex++ {
		row := Row{}
		for colIndex := 0; colIndex < numCols; colIndex++ {
			originIndex := numCols*rowIndex + colIndex
			row.Cells = append(row.Cells, Cell{
				OriginIndex:   []int{originIndex},
				RelativeIndex: originIndex,
				Class:         DefaultCellClass,
				Id:            cellIdString(originIndex, gridIndex),
			})
		}
		rows = append(rows, row)
	}
	return rows
}

func MakeGrids(NumGrids, singleBoxSize int) map[string]Grid {
	grids := map[string]Grid{}
	sideLength := singleBoxSize * singleBoxSize
	fullGridSize := sideLength * sideLength
	for gridIndex := 0; gridIndex < NumGrids; gridIndex++ {
		g := Grid{
			Rows:      makeGrid(singleBoxSize, gridIndex),
			GridIndex: gridIndex,
			BoxSize:   singleBoxSize,
			Tab:       1, // To change later...
		}
		var err error
		var sm SudokuMap
		switch gridIndex {
		case 0: // Original grid.
			g.GridName = DomainGridName
			sm, err = Identity(fullGridSize)
		default: // Shuffled cells.
			g.GridName = fmt.Sprintf("Random map %d", gridIndex)
			sm, err = createMap(fullGridSize, fullGridSize, gridIndex)
		}
		if err != nil {
			fmt.Printf("error creating identity map: %s\n", err)
			return nil
		}
		g.shuffle(sm)
		grids[g.GridName] = g
	}
	return grids
}

func AbsoluteToRelative(index, sideLength int) (int, int) {
	row := index / sideLength
	col := index % sideLength
	return row, col
}

func NeighborCells(relativeIndex, boxSize int, gridName string) map[string]int {
	// I already hate that I am doing this, but we distinguish the
	// cell the user is hovering over, specifically by giving it a value of 1.
	nbs := map[string]int{}
	sideLength := boxSize * boxSize
	rowN, colN := AbsoluteToRelative(relativeIndex, sideLength)
	mainOriginIndices := map[int]bool{}
	neighborOriginIndices := map[int]bool{}
	curGridIndex := Grids[gridName].GridIndex
	// use relative index for grid the user is hovering over and
	// generated origin for other.
	for curAbsIndex := 0; curAbsIndex < sideLength*sideLength; curAbsIndex++ {
		r, c := AbsoluteToRelative(curAbsIndex, sideLength)
		// neighbors have the same row, column, or equivalence class
		// mod size in row AND column
		sameRow := r == rowN
		sameCol := c == colN
		sameBox := r/boxSize == rowN/boxSize && c/boxSize == colN/boxSize

		if sameRow || sameCol || sameBox {
			if r == rowN && c == colN {
				for _, originIndex := range Grids[gridName].Rows[r].Cells[c].OriginIndex {
					mainOriginIndices[originIndex] = true
				}
				nbs[cellIdString(curAbsIndex, curGridIndex)] = 1
			} else {
				for _, originIndex := range Grids[gridName].Rows[r].Cells[c].OriginIndex {
					neighborOriginIndices[originIndex] = true
				}
				nbs[cellIdString(curAbsIndex, curGridIndex)] = 0
			}
		}
	}
	for name, g := range Grids {
		if name == gridName {
			continue
		}
		// My never-nester senses are tingling...
		for _, row := range g.Rows {
			for _, cell := range row.Cells {
				for _, originIndex := range cell.OriginIndex {
					if _, ok := mainOriginIndices[originIndex]; ok {
						nbs[cellIdString(cell.RelativeIndex, g.GridIndex)] = 1
						break
					}
					if _, ok := neighborOriginIndices[originIndex]; ok {
						nbs[cellIdString(cell.RelativeIndex, g.GridIndex)] = 0
						// we have to keep searching in case another
						// origin intersects the currently hovered
						// cell.
					}
				}
			}
		}
	}
	return nbs
}

const (
	gridTmpl = "src/templates/grid.html"
	cellTmpl = "src/templates/cell.html"
)

var (
	updateCount = 0
)

// map[string]map[string]Grid -- do you see what you make me do,
// Go templates?
func RenderGrid(w http.ResponseWriter, gridData map[string]map[string]Grid) error {
	fmt.Printf("update #%d\n", updateCount)
	updateCount++
	tmpl, err := template.ParseFiles(gridTmpl, cellTmpl)
	if err != nil {
		return fmt.Errorf("error parsing template: %s", err)
	}
	if err := tmpl.Execute(w, gridData); err != nil {
		return fmt.Errorf("error rendering grids: %s", err)
	}
	return nil
}
