package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cristiano-pacheco/go-api/core/user"
	"github.com/cristiano-pacheco/go-api/web/common"
	"github.com/cristiano-pacheco/go-api/web/middlewares"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// MakeUserHandlers create all user resource handlers
func MakeUserHandlers(r *mux.Router, n *negroni.Negroni, service user.UseCase, jwtKey *jwt.HMACSHA) {
	r.Handle("/v1/users", n.With(
		middlewares.CheckAuthentication(jwtKey),
		negroni.Wrap(getAllUsers(service)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/v1/users/{id}", n.With(
		middlewares.CheckAuthentication(jwtKey),
		negroni.Wrap(getUser(service)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/v1/users", n.With(
		middlewares.CheckAuthentication(jwtKey),
		negroni.Wrap(storeUser(service)),
	)).Methods("POST", "OPTIONS")

	r.Handle("/v1/users/{id}", n.With(
		middlewares.CheckAuthentication(jwtKey),
		negroni.Wrap(updateUser(service)),
	)).Methods("PUT", "OPTIONS")

	r.Handle("/v1/users/{id}", n.With(
		middlewares.CheckAuthentication(jwtKey),
		negroni.Wrap(removeUser(service)),
	)).Methods("DELETE", "OPTIONS")
}

func getAllUsers(service user.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		all, err := service.GetAll()
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(all)
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func getUser(service user.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		u, err := service.Get(id)
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(u)
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func storeUser(service user.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var u user.User

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}
		//@TODO precisamos validar os dados antes de salvar na base de dados. Pergunta: Como fazer isso?
		err = service.Store(&u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}

func updateUser(service user.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		var u user.User

		err = json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		u.ID = id
		err = service.Update(&u)
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func removeUser(service user.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		err = service.Remove(id)
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
