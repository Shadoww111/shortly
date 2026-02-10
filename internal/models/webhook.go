package models

import "time"

type Webhook struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"` // click, link.created, link.expired
	Secret    string    `json:"-"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type WebhookPayload struct {
	Event     string      `json:"event"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type CreateWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}
