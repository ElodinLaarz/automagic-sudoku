package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
)

const (
	defaultCellClass     = "cell"
	highlightedCellClass = "cell-highlighted"

	gridSize = 3
	numGrids = 3
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
	Rows     []Row
	GridSize int
}

func hover(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "yes")
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
				Class:       defaultCellClass,
				Id:          cellIdString(originIndex, gridIndex),
				Value:       0,
			})
		}
		rows = append(rows, row)
	}
	return rows
}

func makeGrids(numGrids, gridSize int) []Grid {
	var grids []Grid
	for i := 0; i < numGrids; i++ {
		g := Grid{
			Rows:     makeGrid(gridSize, i),
			GridSize: gridSize,
		}
		if i != 0 {
			g.shuffleRows(int64(i))
		}
		grids = append(grids, g)
	}
	return grids
}

func renderGrids(w http.ResponseWriter, gridData map[string][]Grid) error {
	tmpl, err := template.ParseFiles("templates/grid.html")
	// fmt.Print("re-render grid")
	if err != nil {
		return err
	}
	fmt.Printf("gridData: %v\n", gridData)
	return tmpl.Execute(w, gridData)
}

func absoluteToRelative(index, sideLength int) (int, int) {
	row := index / sideLength
	col := index % sideLength
	return row, col
}

func neighborCells(index, size int) map[string]bool {
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
			for curGrid := 0; curGrid < numGrids; curGrid++ {
				nbs[cellIdString(curAbsIndex, curGrid)] = true
			}
		}
	}
	return nbs
}

func highlightHandler(w http.ResponseWriter, r *http.Request) {
	originIndex, err := strconv.Atoi(r.URL.Query().Get("originIndex"))
	if err != nil {
		fmt.Fprintf(w, "Error parsing id: %s", err)
		return
	}
	neighbors := neighborCells(originIndex, gridSize)

	data := makeGrids(numGrids, gridSize)
	for _, grid := range data {
		for rowIndex, row := range grid.Rows {
			for colIndex, cell := range row.Cells {
				if _, ok := neighbors[cell.Id]; ok {
					grid.Rows[rowIndex].Cells[colIndex].Class = highlightedCellClass
					grid.Rows[rowIndex].Cells[colIndex].Value = 1
				}
			}
		}
	}
	gridData := make(map[string][]Grid)
	gridData["grid"] = data
	if err := renderGrids(w, gridData); err != nil {
		fmt.Fprintf(w, "Error highlighting id %d: %s", originIndex, err)
	}
}

func gridHandler(w http.ResponseWriter, r *http.Request) {
	gridData := map[string][]Grid{
		"grid": makeGrids(numGrids, gridSize),
	}
	if err := renderGrids(w, gridData); err != nil {
		fmt.Fprintf(w, "Error rendering grids: %s", err)
	}
}

func main() {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println("ParseFiles:", err)
		return
	}
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./"))
	router.Handle("/css/", fileServer) // Serve CSS files

	router.HandleFunc("/grids", gridHandler)
	router.HandleFunc("/hover", hover)
	router.HandleFunc("/highlight", highlightHandler)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		grids := map[string][]Grid{
			"grid": makeGrids(numGrids, gridSize),
		}
		err := tmpl.Execute(w, grids)
		if err != nil {
			fmt.Println("Execute:", err)
		}
	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
