package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"

	borrowingsvc "library_system/internal/application/borrowing"
	"library_system/internal/domain/borrowing"
)

type BorrowHandler struct {
	svc *borrowingsvc.Service
}

func NewBorrowHandler(svc *borrowingsvc.Service) *BorrowHandler {
	return &BorrowHandler{svc: svc}
}

func (h *BorrowHandler) BorrowBook(w http.ResponseWriter, r *http.Request) {
	var input borrowingsvc.BorrowBookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.svc.BorrowBook(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, borrowing.ErrBookNotAvailable):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, borrowing.ErrMemberSuspended):
			writeError(w, http.StatusForbidden, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *BorrowHandler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	out, err := h.svc.ReturnBook(r.Context(), borrowingsvc.ReturnBookInput{BorrowID: id})
	if err != nil {
		if errors.Is(err, borrowing.ErrNotFound) {
			writeError(w, http.StatusNotFound, "borrow record not found")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, out)
}