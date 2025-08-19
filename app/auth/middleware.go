package auth

import (
	"context"
	"gochess/lib/jwt"
	"log"
	"net/http"
)

const kClaimsKey = "ctxClaims"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(kAuthCookie)
		if err != nil {
			log.Printf("Failed to extract cookie from request: %v\n", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if cookie.Value == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		var claims UserClaims
		if err := jwt.ParseToken(cookie.Value, []byte(kJwtSecret), &claims); err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), kClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Claims(ctx context.Context) (UserClaims, bool) {
	c, ok := ctx.Value(kClaimsKey).(UserClaims)
	return c, ok
}
