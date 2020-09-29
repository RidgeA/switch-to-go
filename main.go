package main

import (
	"github.com/RidgeA/switch-to-go/cache"
	"github.com/RidgeA/switch-to-go/handlers"
	"github.com/RidgeA/switch-to-go/middlewares"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

var port string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func main() {

	logger := middlewares.NewLogger(log.New(os.Stdout, "[s2go]", 0))

	//todo: move to config
	duration, _ := time.ParseDuration("5s")
	storage := cache.NewMemoryCache(duration)
	cacher := middlewares.NewCacher(storage)

	r := mux.NewRouter()
	r.HandleFunc("/", logger(handlers.Index))
	r.HandleFunc("/materials", logger(cacher(handlers.Materials)))

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
