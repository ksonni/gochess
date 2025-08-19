package auth

import (
	"encoding/json"
	"gochess/lib/env"
	"gochess/lib/jwt"
	"log"
	"net/http"
	"time"
)

const kAuthCookie = "auth"

var kJwtSecret = env.MustEnv("JWT_SECRET")
var kCookieExpiryDays = 365

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	c := NewBasicClaims()

	token, err := jwt.CreateToken(c, []byte(kJwtSecret))
	if err != nil {
		log.Printf("Failed to create new user token: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	maxAge := 24 * time.Hour * 365
	expires := time.Now().Add(maxAge)

	cookie := &http.Cookie{
		Name:     kAuthCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  expires,
	}

	log.Printf("Registered user %s\n", c.Id)

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := Claims(r.Context())
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	_ = json.NewEncoder(w).Encode(claims)
}
