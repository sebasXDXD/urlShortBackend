package routes

import (
	"net/http"
	"urlShortenerBack/controllers"
	services "urlShortenerBack/services/users"

	"github.com/gorilla/mux"
)

func SetupRoutes(userService services.UserService) *mux.Router {
	r := mux.NewRouter()
	userController := controllers.NewUserController(userService)
	//ruta para traer tareas
	r.HandleFunc("/users", userController.Index).Methods(http.MethodGet)
	// ruta para agregar tareas
	r.HandleFunc("/user", userController.Create).Methods(http.MethodPost)
	//Ruta de login
	r.HandleFunc("/login", userController.Login).Methods(http.MethodPost)
	return r
}
