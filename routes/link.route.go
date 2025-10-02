package routes

import (
	"net/http"
	"urlShortenerBack/auth"
	"urlShortenerBack/controllers"
	linksService "urlShortenerBack/services/links"

	"github.com/gorilla/mux"
)

func setupLinkRoutes(r *mux.Router, linkService linksService.LinkService) {
	linkController := controllers.NewLinkController(linkService)

	// CRUD de links
	createLinkHandler := http.HandlerFunc(linkController.Create)
	r.Handle("/link", auth.AuthMiddleware(createLinkHandler)).Methods(http.MethodPost)
	r.HandleFunc("/links", linkController.Index).Methods(http.MethodGet)
	r.Handle("/user/links", auth.AuthMiddleware(http.HandlerFunc(linkController.GetUserLinks))).Methods(http.MethodGet)
	r.HandleFunc("/link/{link}", linkController.GetLink).Methods(http.MethodGet)
	r.HandleFunc("/{shortCode}", linkController.Redirect).Methods(http.MethodGet)
	r.Handle("/link/id/{id}", auth.AuthMiddleware(http.HandlerFunc(linkController.GetByID))).Methods(http.MethodGet)
	updateLinkHandler := http.HandlerFunc(linkController.Update)
	r.Handle("/link/{id}", auth.AuthMiddleware(updateLinkHandler)).Methods(http.MethodPut)

	// 📊 Rutas de estadísticas
	r.Handle("/links/stats", auth.AuthMiddleware(http.HandlerFunc(linkController.GetStatsByUserID))).Methods(http.MethodGet)
	r.Handle("/links/clicks", auth.AuthMiddleware(http.HandlerFunc(linkController.GetClicksByMonth))).Methods(http.MethodGet)
	r.Handle("/links/top", auth.AuthMiddleware(http.HandlerFunc(linkController.GetTopLinksByUser))).Methods(http.MethodGet)
	r.Handle("/links/created", auth.AuthMiddleware(http.HandlerFunc(linkController.GetLinksCreatedByMonth))).Methods(http.MethodGet)
	r.Handle("/links/recent", auth.AuthMiddleware(http.HandlerFunc(linkController.GetRecentLinksByUser))).Methods(http.MethodGet)
}
