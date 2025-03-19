package routes

import (
	"net/http"
	"urlShortenerBack/controllers"
	linksService "urlShortenerBack/services/links"
	userService "urlShortenerBack/services/users"

	"github.com/gorilla/mux"
)

// SetupRoutes configura y devuelve el enrutador principal de la aplicación.
func SetupRoutes(userService userService.UserService, linkService linksService.LinkService) *mux.Router {
	r := mux.NewRouter()

	// Configuración de rutas para usuarios
	setupUserRoutes(r, userService)

	// Configuración de rutas para enlaces
	setupLinkRoutes(r, linkService)

	return r
}

// setupUserRoutes configura las rutas relacionadas con usuarios.
func setupUserRoutes(r *mux.Router, userService userService.UserService) {
	userController := controllers.NewUserController(userService)
	r.HandleFunc("/users", userController.Index).Methods(http.MethodGet)
	r.HandleFunc("/user", userController.Create).Methods(http.MethodPost)
	r.HandleFunc("/googleUser", userController.CreateGoogleUser).Methods(http.MethodPost)
	r.HandleFunc("/login", userController.Login).Methods(http.MethodPost)
	r.HandleFunc("/loginGoogle", userController.LoginGoogle).Methods(http.MethodPost)
}
