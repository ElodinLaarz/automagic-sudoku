package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

    spb "github.com/ElodinLaarz/automagic-sudoku/src/proto/sudoku"
)

func generateSudoku() [9][9]int {
    return nil
}

func validSudoku(s *sudokuGrid) bool {
    return false
}

type sudokuGrid [9][9]int

func (s *sudokuGrid) String() string {
    prettyGrid := ""
    for i := 0; i < 9; i++ {
        for j := 0; j < 9; j++ {
            prettyGrid += fmt.Sprintf("%d", s[i][j])
        }
    }
    return prettyGrid
}

func (s *sudokuGrid) rows() [9][9]int {
    return *s
}

// This is so inefficient if we ever do it more than once...
func (s *sudokuGrid) cols() [9][9]int {
    var cols [9][9]int
    for i := 0; i < 9; i++ {
        for j := 0; j < 9; j++ {
            cols[j][i] = s[i][j]
        }
    }
    return cols
}

type pairedGrids struct {
    Grid1 sudokuGrid
    Grid2 sudokuGrid 
}

func main() {
	h1 := func(w http.ResponseWriter, r *http.Request) {
        data := pairedGrids{
            grid1: sudokuGrid{generateSudoku()},
            grid2: sudokuGrid{generateSudoku()},
        }
		tmpl := template.Must(template.ParseFiles("index.html"))
		_ = tmpl.Execute(w, nilListener )
	}
	http.ListenAndServe(":8080", nil)
}
