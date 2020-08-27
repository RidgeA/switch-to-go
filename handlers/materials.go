package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

var materialsTemplate = template.Must(template.New("materials.gohtml").ParseFiles("./views/materials.gohtml"))

func Materials(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("./data/materials-second-module.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "can not read file: %s", err)
	}

	var links []string
	if err := json.Unmarshal(content, &links); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "can not parse file: %s", err)
	}

	if err := materialsTemplate.Execute(w, links); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "error: %s", err)
	}
}
