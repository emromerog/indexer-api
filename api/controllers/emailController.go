package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/emromerog/indexer-api/pkg/utils"
	"github.com/emromerog/indexer-api/pkg/zincsearch"
	"github.com/go-chi/chi/v5"
)

/*Gets all emails without search filtering*/
func GetAllEmails(w http.ResponseWriter, r *http.Request) {
	emails, err := zincsearch.SearchData("", "alldocuments")
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
func SearchBookByTerm(w http.ResponseWriter, r *http.Request) {
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
