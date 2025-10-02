package main

import (
	"fmt"
	"net/http"
	"os"
	"urlShortenerBack/db"
	"urlShortenerBack/repositories"
	"urlShortenerBack/routes"
	serviceslinks "urlShortenerBack/services/links"
	services "urlShortenerBack/services/users"

	"github.com/gorilla/handlers"
)

func main() {

	dbInstance, err := db.ConnectDb()

	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return
	}

	linkService := serviceslinks.NewLinkService(repositories.NewLinkRepository(dbInstance))
	userService := services.NewUserService(repositories.NewUserRepository(dbInstance))

	// Configurar el enrutador (router)
	router := routes.SetupRoutes(userService, linkService)

	// Configurar los encabezados CORS usando gorilla/handlers
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{
		"http://localhost:3000",
		"https://urlshortfrontend.onrender.com",
	})

	// Utilizar el middleware para manejar CORS
	handler := handlers.CORS(headers, methods, origins)(router)

	// Leer el puerto de la variable de entorno o usar 8000 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Servidor escuchando en el puerto " + port + "...")

	// Iniciar el servidor
	err = http.ListenAndServe(":"+port, handler)

	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
