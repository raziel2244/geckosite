package handlers

import (
	"context"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/raziel2244/geckosite/database"
	"github.com/raziel2244/geckosite/database/model"
	"github.com/raziel2244/geckosite/s3"
)

// Home returns the home page.
func Home(w http.ResponseWriter, r *http.Request) {
	type card struct{ Title, Path, Image string }
	type page struct {
		Cards []card
	}

	data := page{
		Cards: []card{
			{"Crested Geckos", "/geckos/crested", ""},
			{"Gargoyle Geckos", "/geckos/gargoyle", ""},
			{"Leopard Geckos", "/geckos/leopard", ""},
		},
	}

	rand.Seed(time.Now().Unix())

	for c, card := range data.Cards {
		splitPath := strings.Split(card.Path, "/")

		var species model.Species
		database.DB.
			Where(&model.Species{Order: splitPath[1], Type: splitPath[2]}).
			First(&species)

		var animals []*model.Animal
		// check for sale
		// check holdbacks
		// check personal

		wheres := []string{
			"status = 'For Sale'",
			"status = 'Holdback'",
			"status IN ('Non-Breeder','Breeder','Future Breeder')",
		}

		for _, where := range wheres {
			database.DB.Where("species_id = ? AND "+where, species.ID.String()).Find(&animals)
			if len(animals) > 0 {
				break
			}
		}

		if len(animals) == 0 {
			data.Cards[c].Image = "/static/img/coming-soon.jpg"
			continue
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(animals), func(i, j int) {
			animals[i], animals[j] = animals[j], animals[i]
		})

		for _, animal := range animals {
			// get images for animal
			var images []string
			ch := s3.Client.ListObjects(
				context.Background(),
				splitPath[1],
				minio.ListObjectsOptions{
					Prefix:    splitPath[2] + "/" + animal.ID.String(),
					Recursive: true,
				},
			)
			for object := range ch {
				path := "/s3/" + splitPath[1] + "/" + object.Key
				images = append(images, path)
			}

			// pick random one for card
			if len(images) > 0 {
				data.Cards[c].Image = images[0]
				break
			}
		}
	}

	lp, hp := "templates/layout.gohtml", "templates/home.gohtml"
	tmpl := template.Must(template.ParseFiles(lp, hp))
	tmpl.ExecuteTemplate(w, "layout", data)
}
