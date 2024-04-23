package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/emromerog/indexer-api/pkg/models"
	"github.com/emromerog/indexer-api/pkg/utils"
	"github.com/emromerog/indexer-api/pkg/zincsearch"
	"github.com/go-chi/chi/v5"
)

/*Gets all emails without search filtering*/
func GetAllEmails(w http.ResponseWriter, r *http.Request) {
	emails, err := zincsearch.SearchData("", "match_all")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al buscar emails: %v", err), http.StatusInternalServerError)
		return
	}

	// Serializar la respuesta en formato JSON
	jsonResponse, err := utils.ConvertToJson(emails)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al serializar la respuesta: %v", err), http.StatusInternalServerError)
		return
	}

	// Establecer encabezados y escribir la respuesta en el cuerpo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

/*Gets all emails according to the search parameter*/
func SearchEmailsByTerm(w http.ResponseWriter, r *http.Request) {
	termParam := chi.URLParam(r, "term")
	cleanedTerm := strings.TrimSpace(termParam)

	if termParam == "" {
		http.Error(w, "A 'term' parameter is required", http.StatusBadRequest)
		return
	}

	emails, err := zincsearch.SearchData(cleanedTerm, "match")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al buscar emails: %v", err), http.StatusInternalServerError)
		return
	}

	// Serializar la respuesta en formato JSON
	jsonResponse, err := utils.ConvertToJson(emails)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al serializar la respuesta: %v", err), http.StatusInternalServerError)
		return
	}

	// Establecer encabezados y escribir la respuesta en el cuerpo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

/*Gets all emails according to the search parameter by JSON*/
func SearchEmailsByJSON(w http.ResponseWriter, r *http.Request) {
	// Decodificar el JSON del cuerpo de la solicitud en la estructura SearchRequest
	fmt.Println("controlador")
	var req models.EmailSearchRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al decodificar JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validar los campos requeridos en la solicitud
	if req.Term == "" {
		http.Error(w, "El campo 'term' es obligatorio", http.StatusBadRequest)
		return
	}

	emails, err := zincsearch.SearchDataJSON(req.Term, req.Page, req.ItemsPerPage, "match")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al buscar emails: %v", err), http.StatusInternalServerError)
		return
	}

	// Serializar la respuesta en formato JSON
	jsonResponse, err := utils.ConvertToJson(emails)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al serializar la respuesta: %v", err), http.StatusInternalServerError)
		return
	}

	// Establecer encabezados y escribir la respuesta en el cuerpo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
