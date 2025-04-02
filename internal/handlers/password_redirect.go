package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/shortly/internal/services"
)

type PasswordRequest struct {
	Password string `json:"password"`
}

// PasswordRedirect handles password-protected link access.
// POST /{code}/unlock with {"password": "..."}
func (h *LinkHandler) PasswordRedirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	var req PasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request", http.StatusBadRequest)
		return
	}

	// get link with password hash
	var passwordHash string
	var originalURL string
	var linkID int
	err := h.links.DB().QueryRow(r.Context(),
		"SELECT id, original_url, password_hash FROM links WHERE short_code=$1 AND is_active=true",
		code,
	).Scan(&linkID, &originalURL, &passwordHash)
	if err != nil {
		writeError(w, "not found", http.StatusNotFound)
		return
	}

	if passwordHash == "" {
		// no password needed, just redirect
		writeJSON(w, map[string]string{"url": originalURL}, http.StatusOK)
		return
	}

	if !services.VerifyLinkPassword(req.Password, passwordHash) {
		writeError(w, "wrong password", http.StatusForbidden)
		return
	}

	// record click
	go h.clicks.Record(r.Context(), linkID, r.RemoteAddr, r.UserAgent(), r.Referer())

	writeJSON(w, map[string]string{"url": originalURL}, http.StatusOK)
}
