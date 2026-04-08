package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"

	booksvc "library_system/internal/application/book"
	"library_system/internal/domain/book"
)

type BookHandler struct {
	svc *booksvc.Service
}

func NewBookHandler(svc *booksvc.Service) *BookHandler {
	return &BookHandler{svc: svc}
}

func (h *BookHandler) AddBook(w http.ResponseWriter, r *http.Request) {
	var input booksvc.AddBookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.svc.AddBook(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	out, err := h.svc.GetBook(r.Context(), id)
	if err != nil {
		if errors.Is(err, book.ErrNotFound) {
			writeError(w, http.StatusNotFound, "book not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	writeJSON(w, http.StatusOK, out)
}