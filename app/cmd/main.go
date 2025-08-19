package main

import (
	"fmt"
	"gochess/db"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const kPort = 8080

func main() {
	log.Printf("Running migrations")
	if err := db.Migrate(); err != nil {
		log.Fatalf("DB migration failed: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "echo")
	})

	port := fmt.Sprintf(":%d", kPort)
	fmt.Printf("Server listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
