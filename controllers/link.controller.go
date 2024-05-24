package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"urlShortenerBack/entities"
	services "urlShortenerBack/services/links"
	"urlShortenerBack/utils"
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

func (c LinkController) Create(w http.ResponseWriter, r *http.Request) {
	// Obtener el "id" del contexto
	userID := r.Context().Value(utils.UserIDKey)

	// Debug: imprimir el valor de userID y su tipo
	fmt.Printf("Valor de userID en el contexto: %v, tipo: %T\n", userID, userID)

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

	// Debug: imprimir el valor convertido de userIDInt
	fmt.Printf("ID de usuario convertido: %d\n", userIDInt)

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdLink)
}
