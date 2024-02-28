package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/emromerog/indexer-api/api/routes"
	"github.com/emromerog/indexer-api/pkg/utils"
)

/*type Server struct {
	Router *chi.Mux
}*/

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	return r
}

func InitializeServer() error {

	r := configureRouter()

	err := startServer(r, utils.DefaultPort)
	if err != nil {
		return fmt.Errorf("error starting the server: %s", err)
	}

	return nil
}

func configureRouter() *chi.Mux {
	// Create a new Chi router
	r := newRouter()

	useMiddleware(r)

	// Mount the routes from the "routes" package under the "/api/v1" path as a subrouter
	r.Mount("/api/v1/emails", routes.EmailRoutes())

	r.Mount("/debug", middleware.Profiler())

	// Define a basic route for the root path
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("¡Welcome!\n"))
		w.Write([]byte("GET -> baseURL/api/v1/emails\n"))
		w.Write([]byte("POST -> baseURL/api/v1/emails/search\n"))
	})

	return r
}

func useMiddleware(r *chi.Mux) {
	// The Logger middleware will print detailed information about the request to the server's console..
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}

func startServer(r *chi.Mux, port string) error {
	fmt.Printf("Server started at http://localhost: %s ", port)

	// Add the main router (r) to the HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting the server: %s", err)
	}

	return nil
}
