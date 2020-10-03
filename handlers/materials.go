package handlers

import (
	"context"
	"fmt"
	"github.com/RidgeA/switch-to-go/db/repository"
	"html/template"
	"net/http"
)

type (
	MaterialsGetter interface {
		GetMaterials(ctx context.Context) ([]repository.Material, error)
	}
)

var materialsHtml = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Switch2Go 2nd module Materials</title>
</head>
<body>
<span>Materials</span>
<ul>
    {{range .}}
    <li>
        <a href="{{.Url}}">{{.Url}}</a>
    </li>
    {{end}}
</ul>
</body>
</html>
`
var materialsTemplate = template.Must(template.New("materials").Parse(materialsHtml))

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
