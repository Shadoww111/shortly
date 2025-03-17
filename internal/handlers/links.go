package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/shortly/internal/middleware"
	"github.com/shortly/internal/models"
	"github.com/shortly/internal/services"
)

type LinkHandler struct {
	links  *services.LinkService
	clicks *services.ClickService
}

func NewLinkHandler(links *services.LinkService, clicks *services.ClickService) *LinkHandler {
	return &LinkHandler{links: links, clicks: clicks}
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req models.CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	link, err := h.links.Create(r.Context(), userID, req)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, link, http.StatusCreated)
}

func (h *LinkHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	resp, err := h.links.ListByUser(r.Context(), userID, page, perPage)
	if err != nil {
		writeError(w, "error fetching links", http.StatusInternalServerError)
		return
	}

	writeJSON(w, resp, http.StatusOK)
}

func (h *LinkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	linkID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.links.Delete(r.Context(), linkID, userID); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, map[string]string{"msg": "deleted"}, http.StatusOK)
}

func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, linkID, err := h.links.Resolve(r.Context(), code)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	// record click async
	go h.clicks.Record(r.Context(), linkID, r.RemoteAddr, r.UserAgent(), r.Referer())

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func (h *LinkHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	linkID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "invalid id", http.StatusBadRequest)
		return
	}
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days < 1 || days > 365 {
		days = 30
	}

	stats, err := h.clicks.GetStats(r.Context(), linkID, days)
	if err != nil {
		writeError(w, "error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, stats, http.StatusOK)
}
