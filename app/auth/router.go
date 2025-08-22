package auth

import (
	"github.com/go-chi/chi/v5"
)

var c = NewController()

func RegisterPublicRoutes(r chi.Router) {
	r.Post("/auth/register", c.registrationHandler)
}

func RegisterRoutes(r chi.Router) {
	r.Get("/auth/me", c.meHandler)
}
