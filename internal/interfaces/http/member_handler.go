package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"

	membersvc "library_system/internal/application/member"
	"library_system/internal/domain/member"
)

type MemberHandler struct {
	svc *membersvc.Service
}

func NewMemberHandler(svc *membersvc.Service) *MemberHandler {
	return &MemberHandler{svc: svc}
}

func (h *MemberHandler) RegisterMember(w http.ResponseWriter, r *http.Request) {
	var input membersvc.RegisterMemberInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.svc.RegisterMember(r.Context(), input)
	if err != nil {
		if errors.Is(err, member.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *MemberHandler) GetMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	out, err := h.svc.GetMember(r.Context(), id)
	if err != nil {
		if errors.Is(err, member.ErrNotFound) {
			writeError(w, http.StatusNotFound, "member not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	writeJSON(w, http.StatusOK, out)
}