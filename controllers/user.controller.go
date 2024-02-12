package controllers

import (
	"encoding/json"
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
	tasks, err := c.UserService.GetTasks()
	if err != nil {
		// Manejar el error si lo hubiera, pero no devolver un error HTTP aquí.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder al cliente con el array de tareas en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (c UserController) Create(w http.ResponseWriter, r *http.Request) {
	// Decodificar los datos del cliente (puede variar según el formato que esperes)
	var newTask entities.Users
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Llamar al método CreateTask() del servicio para agregar la nueva tarea
	createdTask, err := c.UserService.CreateTask(newTask)
	if err != nil {
		// Manejar el error si lo hubiera, pero no devolver un error HTTP aquí.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder al cliente con la tarea recién creada en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}
