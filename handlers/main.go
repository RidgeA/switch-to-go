package handlers

import "net/http"

func Main(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
