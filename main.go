package main

import (
	"fmt"
	"net/http"
	"urlShortenerBack/db"
	"urlShortenerBack/repositories"
	"urlShortenerBack/routes"
	services "urlShortenerBack/services/users"
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

	// Aquí puedes agregar manejadores para diferentes rutas usando router.HandleFunc

	// Iniciar el servidor en el puerto 8000
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
