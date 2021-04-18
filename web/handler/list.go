package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cristiano-pacheco/go-api/core/auth"
	"github.com/cristiano-pacheco/go-api/core/list"
	"github.com/cristiano-pacheco/go-api/web/common"
	"github.com/cristiano-pacheco/go-api/web/middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// MakeListHandlers create all resource handlers
func MakeListHandlers(r *mux.Router, n *negroni.Negroni, service list.UseCase, authService *auth.Service) {
	r.Handle("/v1/lists", n.With(
		middleware.CheckAuthentication(authService),
		negroni.Wrap(getAllLists(service)),
	)).Methods("GET", "OPTIONS").Name(auth.GetAllListsAction)

	r.Handle("/v1/lists/{id}", n.With(
		middleware.CheckAuthentication(authService),
		negroni.Wrap(getList(service)),
	)).Methods("GET", "OPTIONS").Name(auth.GetListAction)

	r.Handle("/v1/lists", n.With(
		middleware.CheckAuthentication(authService),
		negroni.Wrap(storeList(service)),
	)).Methods("POST", "OPTIONS").Name(auth.StoreListAction)

	r.Handle("/v1/lists/{id}", n.With(
		middleware.CheckAuthentication(authService),
		negroni.Wrap(updateList(service)),
	)).Methods("PUT", "OPTIONS").Name(auth.UpdateListAction)

	r.Handle("/v1/lists/{id}", n.With(
		middleware.CheckAuthentication(authService),
		negroni.Wrap(removeList(service)),
	)).Methods("DELETE", "OPTIONS").Name(auth.RemoveListAction)

	// list item routes
	r.Handle("/v1/lists/{id}/items", n.With(
		middleware.CheckAuthentication(authService),
		negroni.Wrap(getAllListItems(service)),
	)).Methods("GET", "OPTIONS").Name(auth.GetListItemAction)

	r.Handle("/v1/lists/{id}/items", n.With(
		middleware.CheckAuthentication(authService),
		negroni.Wrap(storeListItem(service)),
	)).Methods("POST", "OPTIONS").Name(auth.StoreLisItemAction)

	r.Handle("/v1/lists/{id}/items/{itemId}", n.With(
		middleware.CheckAuthentication(authService),
		negroni.Wrap(updateListItem(service)),
	)).Methods("PUT", "OPTIONS").Name(auth.UpdateListItemAction)

	r.Handle("/v1/users/{id}/items/{itemId}", n.With(
		middleware.CheckAuthentication(authService),
		negroni.Wrap(removeListItem(service)),
	)).Methods("DELETE", "OPTIONS").Name(auth.RemoveListItemAction)
}

func getAllLists(service list.UseCase) http.Handler {
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

func getList(service list.UseCase) http.Handler {
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

func storeList(service list.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var l list.List

		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		err = service.Store(&l)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}

func updateList(service list.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		var l list.List

		err = json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		l.ID = id
		err = service.Update(&l)
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func removeList(service list.UseCase) http.Handler {
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

func getAllListItems(service list.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		all, err := service.GetAllItems(id)
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

func storeListItem(service list.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var li list.ListItem

		err := json.NewDecoder(r.Body).Decode(&li)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		li.ListID = id

		err = service.StoreItem(&li)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}

func updateListItem(service list.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		itemId, err := strconv.ParseInt(vars["itemId"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		var li list.ListItem

		err = json.NewDecoder(r.Body).Decode(&li)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		li.ID = itemId
		li.ListID = id
		err = service.UpdateItem(&li)
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func removeListItem(service list.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		itemId, err := strconv.ParseInt(vars["itemId"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(common.FormatJSONError(err.Error()))
			return
		}

		err = service.Remove(itemId)
		if err != nil {
			w.Write(common.FormatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
