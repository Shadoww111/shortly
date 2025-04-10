package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/shortly/internal/models"
)

type WebhookService struct {
	client *http.Client
}

func NewWebhookService() *WebhookService {
	return &WebhookService{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *WebhookService) Dispatch(ctx context.Context, webhook *models.Webhook, event string, data interface{}) {
	payload := models.WebhookPayload{
		Event:     event,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("webhook marshal error: %v", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", webhook.URL, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// sign payload
	if webhook.Secret != "" {
		sig := signPayload(body, webhook.Secret)
		req.Header.Set("X-Webhook-Signature", sig)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("webhook delivery failed for %s: %v", webhook.URL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("webhook %s returned %d", webhook.URL, resp.StatusCode)
	}
}

func signPayload(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// TODO: webhook CRUD routes, retry logic, delivery log
