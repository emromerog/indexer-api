package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetAllEmails(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Obtener todos los emails")
}

func SearchBookByTerm(w http.ResponseWriter, r *http.Request) {
	termParam := chi.URLParam(r, "term")

	if termParam == "" {
		http.Error(w, "Se requiere un parámetro 'term'", http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("Buscar email por término: %s", termParam)
	w.Write([]byte(response))
}
