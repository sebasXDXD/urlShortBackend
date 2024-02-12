package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"urlShortenerBack/entities"
	services "urlShortenerBack/services/users"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{UserService: userService}
}

func (c UserController) Index(w http.ResponseWriter, r *http.Request) {
	// Llamar al método GetTask() del servicio
	users, err := c.UserService.GetTasks()
	if err != nil {
		// Manejar el error si lo hubiera, pero no devolver un error HTTP aquí.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder al cliente con el array de tareas en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c UserController) Create(w http.ResponseWriter, r *http.Request) {
	// Copiar el cuerpo de la solicitud
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)

	// Imprimir el contenido del cuerpo de la solicitud antes de decodificar
	body, err := io.ReadAll(tee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Restaurar el cuerpo de la solicitud para que pueda ser leído nuevamente más adelante
	r.Body = io.NopCloser(&buf)

	// Decodificar los datos del cliente (puede variar según el formato que esperes)
	var newUser entities.Users
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Llamar al método CreateUser() del servicio para agregar nuevo usuario
	createdUser, err := c.UserService.CreateUser(newUser)
	if err != nil {
		// Manejar el error si lo hubiera, pero no devolver un error HTTP aquí.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder al cliente con el usuario recién creado en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (c UserController) Login(w http.ResponseWriter, r *http.Request) {
	// Copiar el cuerpo de la solicitud
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)

	// Imprimir el contenido del cuerpo de la solicitud antes de decodificar
	body, err := io.ReadAll(tee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Restaurar el cuerpo de la solicitud para que pueda ser leído nuevamente más adelante
	r.Body = io.NopCloser(&buf)

	// Verificar si el cuerpo de la solicitud contiene tanto username como password
	var user entities.Users
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validar que tanto username como password están presentes
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Se requieren tanto username como password en la solicitud", http.StatusBadRequest)
		return
	}

	// Llamar al método Login() para ejecutar todo lo referente al login de usuario
	userLoged, err := c.UserService.Login(user)
	if err != nil {
		// Manejar el error si lo hubiera, pero no devolver un error HTTP aquí.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Asignar un token al usuario autenticado
	token, err := c.UserService.AuthService.AssignToken(userLoged.ID, userLoged.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Responder al cliente con el usuario recién creado en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Crear un mapa para combinar el usuario y el token
	response := map[string]interface{}{
		"user":  userLoged,
		"token": token,
	}

	json.NewEncoder(w).Encode(response)
}
