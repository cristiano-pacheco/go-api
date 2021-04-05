package middlewares

import (
	"net/http"
	"strings"

	"github.com/cristiano-pacheco/go-api/web/common"
	"github.com/urfave/negroni"
)

// CheckAuthentication middlware
func CheckAuthentication() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		token := extractTokenFromHeaders(r)
		if token == "" {
			w.Write(common.FormatJSONError("Invalid Credentials"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r)
	})
}

func extractTokenFromHeaders(r *http.Request) string {
	a := r.Header.Get("Authorization")

	if a == "" {
		return ""
	}

	parts := strings.Split(a, "Bearer")

	if len(parts) == 2 {
		return parts[1]
	}

	return ""
}
