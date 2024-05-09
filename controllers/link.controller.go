package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"urlShortenerBack/entities"
	services "urlShortenerBack/services/links"
)

type LinkController struct {
	LinkService services.LinkService
}

func NewLinkController(linkService services.LinkService) LinkController {
	return LinkController{LinkService: linkService}
}

func (c LinkController) Index(w http.ResponseWriter, r *http.Request) {
	// Llamar al método GetTasks() del servicio
	links, err := c.LinkService.GetLinks()
	if err != nil {
		// Manejar el error si lo hubiera, pero no devolver un error HTTP aquí.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder al cliente con el array de enlaces en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

func (c LinkController) Create(w http.ResponseWriter, r *http.Request) {
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

	// Llamar al método CreateUser() del servicio para agregar nuevo enlace
	createdLink, err := c.LinkService.CreateLink(newLink)
	if err != nil {
		// Manejar el error si lo hubiera, pero no devolver un error HTTP aquí.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder al cliente con el enlace recién creado en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdLink)
}
