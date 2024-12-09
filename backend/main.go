package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/public/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
