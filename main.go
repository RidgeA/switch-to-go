package main

import (
	"github.com/RidgeA/switch-to-go/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Main)
	r.HandleFunc("/materials", handlers.Materials)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
