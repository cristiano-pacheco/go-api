package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/cristiano-pacheco/go-api/core/authentication"
	"github.com/cristiano-pacheco/go-api/web/common"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/urfave/negroni"
)

// CheckAuthentication middlware
func CheckAuthentication(jwtHash *jwt.HMACSHA) negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		token := extractTokenFromHeaders(r)
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(common.FormatJSONError("Invalid Credentials"))
			return
		}

		var pl authentication.CustomPayload

		now := time.Now()
		iatValidator := jwt.IssuedAtValidator(now)
		expValidator := jwt.ExpirationTimeValidator(now)
		validatePayload := jwt.ValidatePayload(&pl.Payload, iatValidator, expValidator)

		_, err := jwt.Verify([]byte(token), jwtHash, &pl, validatePayload)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(common.FormatJSONError("Invalid Credentials"))
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
		return strings.TrimSpace(parts[1])
	}

	return ""
}
