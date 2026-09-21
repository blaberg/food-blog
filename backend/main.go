package main

import (
	"log"
	"net/http"
	"os"
)

// A development server for previewing the generated site. It serves public/ at
// the root, which is how the site is served in production - the pages use
// root-relative links, so serving them from a subpath breaks navigation.
func main() {
	const addr = "localhost:8080"
	if _, err := os.Stat("public"); err != nil {
		log.Fatalf("public/ not found, run `go run ./cmd/generator` from the repo root first: %v", err)
	}
	log.Printf("Serving public/ on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir("public"))))
}
