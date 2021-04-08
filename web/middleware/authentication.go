package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cristiano-pacheco/go-api/core/authentication"
	"github.com/cristiano-pacheco/go-api/web/common"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gorilla/mux"
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
		userId, err := getUserIdFromToken(token)
		if err != nil || userId == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(common.FormatJSONError("Unable to parse the token data"))
			return
		}
		routeName := mux.CurrentRoute(r).GetName()

		fmt.Println(userId)
		fmt.Println(routeName)

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

func getUserIdFromToken(token string) (int, error) {
	parts := strings.Split(token, ".")
	part := parts[1]
	payload, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return 0, err
	}
	data, err := json_decode(payload)
	if err != nil {
		return 0, err
	}
	return int(data["user_id"].(float64)), nil
}

func json_decode(data []byte) (map[string]interface{}, error) {
	var dat map[string]interface{}
	err := json.Unmarshal(data, &dat)
	return dat, err
}
