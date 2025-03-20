package routes

import (
	"net/http"
	"urlShortenerBack/auth"
	"urlShortenerBack/controllers"
	linksService "urlShortenerBack/services/links"

	"github.com/gorilla/mux"
)

// setupLinkRoutes configura las rutas relacionadas con enlaces.
func setupLinkRoutes(r *mux.Router, linkService linksService.LinkService) {
	linkController := controllers.NewLinkController(linkService)

	createLinkHandler := http.HandlerFunc(linkController.Create)
	r.Handle("/link", auth.AuthMiddleware(createLinkHandler)).Methods(http.MethodPost)
	r.HandleFunc("/links", linkController.Index).Methods(http.MethodGet)
	r.HandleFunc("/link/{link}", linkController.GetLink).Methods(http.MethodGet)
	r.HandleFunc("/{shortCode}", linkController.Redirect).Methods(http.MethodGet)
	r.Handle("/link/id/{id}", auth.AuthMiddleware(http.HandlerFunc(linkController.GetByID))).Methods(http.MethodGet)
	updateLinkHandler := http.HandlerFunc(linkController.Update)
	r.Handle("/link/{id}", auth.AuthMiddleware(updateLinkHandler)).Methods(http.MethodPut)
}
