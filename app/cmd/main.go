package main

import (
	"fmt"
	"log"
	"net/http"
)

const kPort = 8080

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "echo")
	})

	port := fmt.Sprintf(":%d", kPort)
	fmt.Printf("Server listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
