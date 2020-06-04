package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/niko-cb/covid19datascraper/app/controller/handler"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func init() {

}

// @title Covid19 DataScraper
// @description.markdown
// @BasePath /api/
func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.SetHeader("Cache-Control", "no-store"))
	r.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=2592000"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(c.Handler)

	//Routes
	r.Route(handler.APIPathPrefix, func(r chi.Router) {
		r.Route(handler.ScrapeDataAPIBasePath, handler.Scrape)
		r.Route(handler.DialogflowAPIBasePath, handler.Dialogflow)
	})

	// Choose port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s\n", port)
	}
	log.Printf("Listening on port %s\n", port)

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
