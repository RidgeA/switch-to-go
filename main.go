package main

//go:generate go test -cover -coverprofile=./_coverage/coverage.out ./...
//go:generate go tool cover -html=./_coverage/coverage.out -o ./_coverage/index.html

import (
	"context"
	"github.com/RidgeA/switch-to-go/cache"
	"github.com/RidgeA/switch-to-go/db/repository"
	"github.com/RidgeA/switch-to-go/handlers"
	"github.com/RidgeA/switch-to-go/middlewares"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
	"os"
	"time"
)

var port string
var databaseUrl string

func initOptions() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseUrl = os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL environment variable required")
	}
}

func main() {

	initOptions()

	logger := middlewares.NewLogger(log.New(os.Stdout, "[s2go]", 0))

	//todo: move to config
	duration, _ := time.ParseDuration("5s")

	storage := cache.NewMemoryCache(duration)
	cacher := middlewares.NewCacher(storage)

	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("cannot connect to a DB: %s", err.Error())
	}
	materialsRepository := repository.NewMaterialRepository(conn)

	r := mux.NewRouter()
	r.HandleFunc("/",
		logger(
			cacher(
				handlers.Index,
			),
		),
	)
	r.HandleFunc("/materials",
		logger(
			cacher(
				handlers.Materials(materialsRepository),
			),
		),
	)

	r.Handle("/_coverage/", http.StripPrefix("/_coverage/", http.FileServer(http.Dir("_coverage"))))

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
