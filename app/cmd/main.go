package main

import (
	"fmt"
	"gochess/auth"
	"gochess/db"
	"gochess/game"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const kPort = 8080

func main() {
	log.Printf("Running migrations")
	if err := db.Migrate(); err != nil {
		log.Fatalf("DB migration failed: %v", err)
	}

	router := chi.NewRouter()

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(JSONMiddleware)
		r.Group(func(p chi.Router) {
			p.Use(auth.Authenticate)
			auth.RegisterRoutes(p)
			game.RegisterRoutes(p)
		})
		auth.RegisterPublicRoutes(r)
	})

	port := fmt.Sprintf(":%d", kPort)
	fmt.Printf("Server listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
