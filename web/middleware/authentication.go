package middleware

import (
	"encoding/base64"
	"encoding/json"
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
func CheckAuthentication(s *authentication.Service) negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		token := extractTokenFromHeaders(r)
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(common.FormatJSONError("Not Authorized"))
			return
		}

		var pl authentication.CustomPayload

		now := time.Now()
		iatValidator := jwt.IssuedAtValidator(now)
		expValidator := jwt.ExpirationTimeValidator(now)
		validatePayload := jwt.ValidatePayload(&pl.Payload, iatValidator, expValidator)

		_, err := jwt.Verify([]byte(token), s.JWTHash, &pl, validatePayload)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(common.FormatJSONError("Not Authorized"))
			return
		}
		userId, err := getUserIdFromToken(token)
		if err != nil || userId == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(common.FormatJSONError("Unable to parse the token data"))
			return
		}
		routeName := mux.CurrentRoute(r).GetName()

		hasAccess, err := s.HasAccess(userId, routeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		if !hasAccess {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(common.FormatJSONError("Not Authorized"))
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
