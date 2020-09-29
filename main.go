package main

import (
	"context"
	"github.com/RidgeA/switch-to-go/cache"
	"github.com/RidgeA/switch-to-go/db/models"
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

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseUrl = os.Getenv("DATABASE_URL")
}

func main() {

	logger := middlewares.NewLogger(log.New(os.Stdout, "[s2go]", 0))

	//todo: move to config
	duration, _ := time.ParseDuration("5s")

	storage := cache.NewMemoryCache(duration)
	cacher := middlewares.NewCacher(storage)

	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("cannot connect to a DB: %s", err.Error())
	}
	materialsRepository := models.NewMaterialRepository(conn)

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

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
