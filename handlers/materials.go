package handlers

import (
	"context"
	"fmt"
	"github.com/RidgeA/switch-to-go/db/models"
	"html/template"
	"net/http"
)

type (
	MaterialsGetter interface {
		GetMaterials(ctx context.Context) ([]models.Material, error)
	}
)

var materialsTemplate = template.Must(template.New("materials.gohtml").ParseFiles("./views/materials.gohtml"))

func Materials(materials MaterialsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		links, err := materials.GetMaterials(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "can not get list of materials: %s", err)
			return
		}

		if err := materialsTemplate.Execute(w, links); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "error: %s", err)
			return
		}
	}
}
