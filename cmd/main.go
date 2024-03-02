package main

import (
	"fmt"
	"net/http"

	"automagic-sudoku/src/handler"
)

func main() {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("static/"))
	router.Handle("/css/", fileServer) // Serve CSS files

	router.HandleFunc("/", handler.IndexHandler)
	router.HandleFunc("/grids", handler.GridHandler)
	router.HandleFunc("/neighbor", handler.NeighborHandler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
