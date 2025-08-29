package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"urlShortenerBack/entities"
	services "urlShortenerBack/services/links"
	"urlShortenerBack/utils"

	"github.com/gorilla/mux"
)

type LinkController struct {
	LinkService services.LinkService
}

func NewLinkController(linkService services.LinkService) LinkController {
	return LinkController{LinkService: linkService}
}

func (c LinkController) Index(w http.ResponseWriter, r *http.Request) {
	links, err := c.LinkService.GetLinks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}
func (c *LinkController) GetLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkID := vars["link"]

	link, err := c.LinkService.GetLinkByString(linkID)
	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}
func (c *LinkController) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	// 📌 Imprimir shortCode en consola
	fmt.Printf("[%s] ShortCode recibido: %s\n", time.Now().Format(time.RFC3339), shortCode)
	fmt.Printf("ShortCode recibido: %s\n", shortCode)

	// Busca la URL original
	link, err := c.LinkService.GetLinkByString(shortCode)
	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	// 📌 Obtener datos del visitante
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}
	userAgent := r.UserAgent()

	// 📌 Registrar el click (no detenemos la redirección si falla)
	if err := c.LinkService.RegisterClick(link.ID, ip, userAgent); err != nil {
		fmt.Printf("Error registrando click: %v\n", err)
	}

	// 📌 Redirigir a la URL original
	http.Redirect(w, r, link.RedirectTo, http.StatusFound)
}

func (c LinkController) Create(w http.ResponseWriter, r *http.Request) {
	// Obtener el "id" del contexto
	userID := r.Context().Value(utils.UserIDKey)

	if userID == nil {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	// Asegúrate de que userID es del tipo correcto (float64)
	userIDFloat, ok := userID.(float64)
	if !ok {
		fmt.Println("Tipo de usuario no válido en el contexto, se esperaba float64")
		http.Error(w, "Tipo de usuario no válido en el contexto", http.StatusUnauthorized)
		return
	}
	userIDInt := int(userIDFloat)

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
	var newLink entities.Link
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&newLink); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Asignar el userID extraído al nuevo enlace
	newLink.UserCreatedID = userIDInt

	// Llamar al método CreateUser() del servicio para agregar nuevo enlace
	createdLink, err := c.LinkService.CreateLink(newLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder al cliente con el enlace recién creado en formato JSON
	response := map[string]interface{}{
		"message": "Enlace creado exitosamente",
		"link":    createdLink,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c LinkController) GetByID(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del enlace desde los parámetros de la URL
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID a entero
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Obtener el enlace desde el repositorio
	link, err := c.LinkService.GetLinkByID(id)
	if err != nil {
		http.Error(w, "Error al obtener el enlace", http.StatusInternalServerError)
		return
	}

	if link.ID == 0 {
		http.Error(w, "Enlace no encontrado", http.StatusNotFound)
		return
	}

	// Responder con el enlace en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

func (c LinkController) Update(w http.ResponseWriter, r *http.Request) {
	// Obtener el "id" del contexto
	userID := r.Context().Value(utils.UserIDKey)

	if userID == nil {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	// Asegurar que userID es del tipo correcto (float64)
	userIDFloat, ok := userID.(float64)
	if !ok {
		http.Error(w, "Tipo de usuario no válido", http.StatusUnauthorized)
		return
	}
	userIDInt := int(userIDFloat)

	// Obtener el ID del enlace desde los parámetros de la URL
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID a entero
	linkID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Decodificar el cuerpo de la solicitud
	var updateData struct {
		Name       string `json:"name"`
		RedirectTo string `json:"redirect_to"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	// Validar que los campos no estén vacíos
	if updateData.Name == "" || updateData.RedirectTo == "" {
		http.Error(w, "Nombre y URL son obligatorios", http.StatusBadRequest)
		return
	}

	// Obtener el enlace actual
	link, err := c.LinkService.GetLinkByID(linkID)
	if err != nil {
		http.Error(w, "Error al obtener el enlace", http.StatusInternalServerError)
		return
	}

	if link.ID == 0 {
		http.Error(w, "Enlace no encontrado", http.StatusNotFound)
		return
	}

	// Verificar que el usuario autenticado es el creador del enlace
	if link.UserCreatedID != userIDInt {
		http.Error(w, "No tienes permisos para modificar este enlace", http.StatusForbidden)
		return
	}

	// Actualizar el enlace
	err = c.LinkService.UpdateLink(linkID, updateData.Name, updateData.RedirectTo)
	if err != nil {
		http.Error(w, "Error al actualizar el enlace", http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Enlace actualizado correctamente"})
}

// Obtener estadísticas generales de un usuario
func (c LinkController) GetStatsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey)
	if userID == nil {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}
	userIDFloat, ok := userID.(float64)
	if !ok {
		http.Error(w, "Tipo de usuario no válido", http.StatusUnauthorized)
		return
	}
	userIDInt := int(userIDFloat)

	stats, err := c.LinkService.GetStatsByUserID(userIDInt)
	if err != nil {
		http.Error(w, "Error obteniendo estadísticas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Obtener clicks por mes (últimos N meses)
func (c LinkController) GetClicksByMonth(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey)
	if userID == nil {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}
	userIDInt := int(userID.(float64))

	monthsStr := r.URL.Query().Get("months")
	months, err := strconv.Atoi(monthsStr)
	if err != nil || months <= 0 {
		months = 6 // valor por defecto
	}

	clicks, err := c.LinkService.GetClicksByMonth(userIDInt, months)
	if err != nil {
		http.Error(w, "Error obteniendo clicks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clicks)
}

// Top N enlaces más clickeados
func (c LinkController) GetTopLinksByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey)
	if userID == nil {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}
	userIDInt := int(userID.(float64))

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5 // valor por defecto
	}

	topLinks, err := c.LinkService.GetTopLinksByUser(userIDInt, limit)
	if err != nil {
		http.Error(w, "Error obteniendo top links", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topLinks)
}

// Links creados por mes (últimos N meses)
func (c LinkController) GetLinksCreatedByMonth(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey)
	if userID == nil {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}
	userIDInt := int(userID.(float64))

	monthsStr := r.URL.Query().Get("months")
	months, err := strconv.Atoi(monthsStr)
	if err != nil || months <= 0 {
		months = 6
	}

	links, err := c.LinkService.GetLinksCreatedByMonth(userIDInt, months)
	if err != nil {
		http.Error(w, "Error obteniendo links creados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

// Links recientes (últimos N creados por el usuario)
func (c LinkController) GetRecentLinksByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey)
	if userID == nil {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}
	userIDInt := int(userID.(float64))

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	recentLinks, err := c.LinkService.GetRecentLinksByUser(userIDInt, limit)
	if err != nil {
		http.Error(w, "Error obteniendo links recientes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recentLinks)
}
