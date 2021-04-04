package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cristiano-pacheco/go-api/core/auth"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// MakeAuthHandlers create all user resource handlers
func MakeAuthHandlers(r *mux.Router, n *negroni.Negroni, service auth.UseCase) {
	r.Handle("/v1/auth", n.With(
		negroni.Wrap(IssueToken(service)),
	)).Methods("POST", "OPTIONS")
}

// IssueToken handler
func IssueToken(service auth.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var ar authRequest

		err := json.NewDecoder(r.Body).Decode(&ar)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(formatJSONError(err.Error()))
			return
		}
		token, err := service.IssueToken(ar.Email, ar.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(formatJSONError(err.Error()))
			return
		}

		err = json.NewEncoder(w).Encode(token)
		if err != nil {
			w.Write(formatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
