package http

import (
	"fmt"
	"net/http"

	//"github.com/emromerog/indexer-api/app/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/emromerog/indexer-api/api/routes"
)

type Server struct {
	Router *chi.Mux
}

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	return r
}

func InitializeServer() error {

	r := configureRouter()

	err := startServer(r, "8080")
	if err != nil {
		return fmt.Errorf("error starting the server: %s", err)
	}

	return nil

}

func useMiddleware(r *chi.Mux) {
	// The Logger middleware will print detailed information about the request to the server's console..
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}

func configureRouter() *chi.Mux {
	// Create a new Chi router
	r := newRouter()

	useMiddleware(r)

	// Mount the routes from the "routes" package under the "/api/v1" path as a subrouter
	r.Mount("/api/v1/emails", routes.EmailRoutes())

	// Define a basic route for the root path
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("¡Welcome!"))
	})

	return r
}

func startServer(r *chi.Mux, port string) error {
	fmt.Printf("Server started at http://localhost:%s\n", port)

	// Add the main router (r) to the HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Retorna el error si hay algún problema al iniciar el servidor
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting the server: %s", err)
	}

	return nil
}
