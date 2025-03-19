package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	qrcode "github.com/skip2/go-qrcode"

	"github.com/shortly/internal/config"
)

type QRHandler struct {
	cfg *config.Config
}

func NewQRHandler(cfg *config.Config) *QRHandler {
	return &QRHandler{cfg: cfg}
}

func (h *QRHandler) Generate(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	sizeStr := r.URL.Query().Get("size")
	size := 256
	if s, err := strconv.Atoi(sizeStr); err == nil && s >= 64 && s <= 1024 {
		size = s
	}

	url := fmt.Sprintf("%s/%s", h.cfg.BaseURL, code)

	png, err := qrcode.Encode(url, qrcode.Medium, size)
	if err != nil {
		writeError(w, "failed to generate qr", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(png)
}
