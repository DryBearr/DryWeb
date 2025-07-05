// ===============================================================
// File: main.go
// Description: Serves the given directory passed as argument
// Author: DryBearr
// ===============================================================

package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Check if path argument is provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: command <directory>")
	}
	dir := os.Args[1]

	fs := http.FileServer(http.Dir(dir))

	http.Handle("/", fs)

	log.Println("Serving", dir, "on http://localhost:9696/")
	err := http.ListenAndServe(":9696", nil)
	if err != nil {
		log.Fatal(err)
	}
}
