package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Route struct {
	path    string
	handler func(w http.ResponseWriter, r *http.Request)
	method  string
}

var routes = make([]Route, 0)

func registerRoute(r Route) {
	routes = append(routes, r)
}

var db *sqlx.DB

func main() {
	dbPath, set := os.LookupEnv("SQLITE_DB_PATH")

	if !set {
		dbPath = "../test.db"
		fmt.Println("[Warn] Using test db file. It's recommended to specify database file path using SQLITE_DB_PATH environment variable.")
	}

	db = sqlx.MustConnect("sqlite3", dbPath)

	_, loaderSet := os.LookupEnv("RUN_DB_LOADER")
	if loaderSet {
		RunDataLoader()
	} else {
		mainRouter := mux.NewRouter()
		r := mainRouter.PathPrefix("/api").Subrouter()

		for _, rt := range routes {
			r.HandleFunc(rt.path, rt.handler).Methods(rt.method)
		}

		mainRouter.Use(mux.CORSMethodMiddleware(r))

		http.ListenAndServe(":8080", mainRouter)
	}
}

// Uses a generic GetAll function for the structs, to create a generic Route Handler.
func handleGetGeneric[T any](getFunc func() ([]T, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		stuff, err := getFunc()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not load items.")
		} else {
			err := json.NewEncoder(w).Encode(&stuff)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err.Error())
				fmt.Fprint(w, "Could not load items.")
			}
		}
	}
}
