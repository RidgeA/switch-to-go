package main

import (
	"github.com/RidgeA/switch-to-go/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var port string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Main)
	r.HandleFunc("/materials", handlers.Materials)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
