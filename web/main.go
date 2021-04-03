package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cristiano-pacheco/go-api/core/user"
	"github.com/cristiano-pacheco/go-api/web/handlers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	dsn := flag.String("dsn", "root:root@/go_api?parseTime=true", "MySQL data source name")
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	service := user.NewService(db)
	r := mux.NewRouter()

	n := negroni.New(
		negroni.NewLogger(),
	)

	//handlers
	handlers.MakeUserHandlers(r, n, service)

	http.Handle("/", r)

	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         *addr,
		Handler:      http.DefaultServeMux,
		ErrorLog:     log.New(os.Stderr, "logger: ", log.Lshortfile),
	}

	srv.ErrorLog.Printf("Server started at %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
