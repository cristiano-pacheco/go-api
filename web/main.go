package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cristiano-pacheco/go-api/core/authentication"
	"github.com/cristiano-pacheco/go-api/core/user"
	"github.com/cristiano-pacheco/go-api/web/handlers"
	"github.com/cristiano-pacheco/go-api/web/middlewares"
	"github.com/gbrlsnchs/jwt/v3"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	dsn := flag.String("dsn", "root:root@/go_api?parseTime=true", "MySQL data source name")
	addr := flag.String("addr", ":4000", "HTTP network address")
	jwtkey := flag.String("jwtkey", "jwt-private-key", "JWT Private Key")
	flag.Parse()

	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Hash used to sign and verify the JWT tokens
	jwtHash := jwt.NewHS256([]byte(*jwtkey))

	userService := user.NewService(db, &user.Validator{})
	authService := authentication.NewService(db, &authentication.Validator{}, jwtHash)
	r := mux.NewRouter()

	n := negroni.New(
		negroni.NewLogger(),
	)
	n.Use(middlewares.SetJSONContentType())

	// handlers
	handlers.MakeUserHandlers(r, n, userService, jwtHash)
	handlers.MakeAuthHandlers(r, n, authService)

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
