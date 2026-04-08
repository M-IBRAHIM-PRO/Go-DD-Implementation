package httphandler

import (
	"encoding/json"
	"net/http"
)

func NewRouter(bookH *BookHandler, memberH *MemberHandler, borrowH *BorrowHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /books", bookH.AddBook)
	mux.HandleFunc("GET /books/{id}", bookH.GetBook)

	mux.HandleFunc("POST /members", memberH.RegisterMember)
	mux.HandleFunc("GET /members/{id}", memberH.GetMember)

	mux.HandleFunc("POST /borrowings", borrowH.BorrowBook)
	mux.HandleFunc("PATCH /borrowings/{id}/return", borrowH.ReturnBook)

	return mux
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}