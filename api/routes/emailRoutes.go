package routes

import (
	"github.com/emromerog/indexer-api/api/controllers"
	"github.com/go-chi/chi/v5"
)

func EmailRoutes() *chi.Mux {
	r := chi.NewRouter()

	/*r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("¡Hola desde la API!"))
	})*/

	r.Get("/", controllers.GetAllEmails)

	r.Get("/search/{term}", controllers.SearchEmailsByTerm)

	r.Post("/", controllers.SearchEmailsByJSON)

	return r
}
