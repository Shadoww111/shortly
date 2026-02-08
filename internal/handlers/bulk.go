package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shortly/internal/middleware"
	"github.com/shortly/internal/models"
)

type BulkRequest struct {
	URLs []models.CreateLinkRequest `json:"urls"`
}

type BulkResponse struct {
	Links  []models.Link `json:"links"`
	Errors []BulkError   `json:"errors,omitempty"`
}

type BulkError struct {
	Index int    `json:"index"`
	URL   string `json:"url"`
	Error string `json:"error"`
}

func (h *LinkHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req BulkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.URLs) == 0 {
		writeError(w, "urls array is empty", http.StatusBadRequest)
		return
	}
	if len(req.URLs) > 50 {
		writeError(w, "max 50 urls per batch", http.StatusBadRequest)
		return
	}

	resp := BulkResponse{}
	for i, item := range req.URLs {
		link, err := h.links.Create(r.Context(), userID, item)
		if err != nil {
			resp.Errors = append(resp.Errors, BulkError{Index: i, URL: item.URL, Error: err.Error()})
			continue
		}
		resp.Links = append(resp.Links, *link)
	}

	writeJSON(w, resp, http.StatusCreated)
}
