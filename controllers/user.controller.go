package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"urlShortenerBack/dtos"
	"urlShortenerBack/entities"
	linksServices "urlShortenerBack/services/links"
	services "urlShortenerBack/services/users"
)

type UserController struct {
	UserService services.UserService
	LinkService linksServices.LinkService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{UserService: userService}
}

func (c UserController) Index(w http.ResponseWriter, r *http.Request) {
	users, err := c.UserService.GetUsers()
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
func (c UserController) CreateGoogleUser(w http.ResponseWriter, r *http.Request) {
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

	// Obtener el token de Google de la solicitud (puede estar en los encabezados o en el cuerpo de la solicitud)
	// googleIDToken := r.Header.Get("Google-ID-Token")
	// fmt.Printf("Google ID Token: %s\n", googleIDToken)
	// if googleIDToken == "" {
	// 	http.Error(w, "Falta el token de Google", http.StatusBadRequest)
	// 	return
	// }

	// Llamar al método CreateGoogleUser() del servicio para agregar el nuevo usuario
	createdUser, err := c.UserService.CreateGoogleUser(r.Context(), newUser)
	if err != nil {
		// Manejar el error si lo hubiera
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
func (c UserController) LoginGoogle(w http.ResponseWriter, r *http.Request) {
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

	// Verificar si el cuerpo de la solicitud contiene tanto username como google_id
	var user entities.Users
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validar que tanto username como google_id están presentes
	if user.Username == "" || user.GoogleID == "" {
		http.Error(w, "Se requieren tanto username como google_id en la solicitud", http.StatusBadRequest)
		return
	}

	// Llamar al método LoginGoogle() para ejecutar todo lo referente al login de usuario con Google
	userLoged, err := c.UserService.LoginGoogle(user)
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

	// Responder al cliente con el usuario recién autenticado en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Crear un mapa para combinar el usuario y el token
	response := map[string]interface{}{
		"user":  userLoged,
		"token": token,
	}

	json.NewEncoder(w).Encode(response)
}

func (uc *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	// 🔐 Extraer usuario desde el token o sesión (esto depende de tu auth)
	userID, err := uc.UserService.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Obtener datos del usuario
	user, err := uc.UserService.GetByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	stats, err := uc.LinkService.GetStatsByUserID(userID)
	if err != nil {
		http.Error(w, "Error getting stats", http.StatusInternalServerError)
		return
	}

	// Estructura combinada
	response := dtos.UserProfileResponse{
		FullName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Email:    user.Email,
		// Company:  user.Company,
		// Country:  user.Country,
		// Phone:    user.Phone,
		// Plan:     user.Plan,
		JoinDate: user.CreatedAt,
		Stats: dtos.StatsData{
			URLsCreated: stats.TotalCreated,
			// URLsLimit:        user.URLLimit,
			TotalClicks:      stats.TotalClicks,
			PopularURL:       stats.MostClickedURL,
			PopularURLClicks: stats.MostClickedCount,
			LastAccess:       stats.LastAccess.Format("02 Jan 2006 15:04"),
		},
	}

	// Enviar respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
