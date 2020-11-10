package handlers

import (
	"context"
	"html/template"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/raziel2244/geckosite/database"
	"github.com/raziel2244/geckosite/database/model"
	"github.com/raziel2244/geckosite/s3"
	"gorm.io/gorm/clause"
)

// Animal returns a page for a single animal.
func Animal(w http.ResponseWriter, r *http.Request) {
	pageData := struct {
		Title, Path string
		Animal      *model.Animal
	}{
		Path: r.URL.Path,
	}

	var animal model.Animal
	database.DB.Model(&model.Animal{}).
		Preload(clause.Associations).
		Where("id = ?", mux.Vars(r)["id"]).
		First(&animal)

	if reflect.ValueOf(animal).IsZero() {
		NotFound(w, r)
		return
	}

	pageData.Animal = &animal

	if len(animal.Name) > 0 {
		pageData.Title = animal.Name
	} else if len(animal.Reference) > 0 {
		pageData.Title = animal.Reference
	} else {
		pageData.Title = animal.Species.Name
	}

	ch := s3.Client.ListObjects(
		context.Background(),
		animal.Species.Order,
		minio.ListObjectsOptions{
			Prefix:    animal.Species.Type + "/" + animal.ID.String(),
			Recursive: true,
		},
	)
	for object := range ch {
		path := "/s3/" + animal.Species.Order + "/" + object.Key
		animal.Images = append(animal.Images, path)
	}

	lp, hp := "templates/layout.gohtml", "templates/animal.gohtml"
	tmpl := template.New("Animal").Funcs(template.FuncMap{
		"formatDate": func(t *time.Time) string {
			return t.Format("02/01/2006")
		},
	})
	tmpl = template.Must(tmpl.ParseFiles(lp, hp))
	tmpl.ExecuteTemplate(w, "layout", pageData)
}
