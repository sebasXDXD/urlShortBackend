package main

import (
	"fmt"
	"net/http"
	"urlShortenerBack/db"
	"urlShortenerBack/repositories"
	"urlShortenerBack/routes"
	services "urlShortenerBack/services/users"

	"github.com/gorilla/handlers"
)

func main() {

	dbInstance, err := db.ConnectDb()

	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return
	}

	userService := services.NewUserService(repositories.NewUserRepository(dbInstance))

	fmt.Println("Servidor escuchando en el puerto 8000...")

	// Configurar el enrutador (router)
	router := routes.SetupRoutes(userService)

	// Configurar los encabezados CORS usando gorilla/handlers
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"}) // Puedes ajustar esto según tus necesidades

	// Utilizar el middleware para manejar CORS
	handler := handlers.CORS(headers, methods, origins)(router)

	// Iniciar el servidor en el puerto 8000
	err = http.ListenAndServe(":8000", handler)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
