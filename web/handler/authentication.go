package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cristiano-pacheco/go-api/core/auth"
	"github.com/cristiano-pacheco/go-api/web/common"
	"github.com/cristiano-pacheco/go-api/web/middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// MakeAuthHandlers create all user resource handlers
func MakeAuthHandlers(r *mux.Router, n *negroni.Negroni, service *auth.Service) {
	r.Handle("/v1/auth", n.With(
		negroni.Wrap(issueToken(service)),
	)).Methods("POST", "OPTIONS")

	r.Handle("/v1/auth/me", n.With(
		middleware.CheckAuthentication(service),
		negroni.Wrap(getAuthenticatedUserData(service)),
	)).Methods("GET", "OPTIONS").Name(auth.UserME)

}

// IssueToken handler
func issueToken(service *auth.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ar authRequest

		err := json.NewDecoder(r.Body).Decode(&ar)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}
		token, err := service.IssueToken(ar.Email, ar.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		err = json.NewEncoder(w).Encode(token)
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

// getAuthenticatedUserData
func getAuthenticatedUserData(service *auth.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("UserID").(int)
		au, err := service.GetUserPermissionsById(userId)
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(au)
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
