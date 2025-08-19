package auth

import (
	"github.com/go-chi/chi/v5"
)

func RegisterPublicRoutes(r chi.Router) {
	r.Post("/auth/register", registrationHandler)
}

func RegisterRoutes(r chi.Router) {
	r.Get("/auth/me", meHandler)
}
