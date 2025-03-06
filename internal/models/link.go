package models

import "time"

type Link struct {
	ID           int        `json:"id"`
	ShortCode    string     `json:"short_code"`
	OriginalURL  string     `json:"original_url"`
	Title        string     `json:"title,omitempty"`
	UserID       int        `json:"user_id"`
	IsActive     bool       `json:"is_active"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	MaxClicks    *int       `json:"max_clicks,omitempty"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ClickCount   int        `json:"click_count,omitempty"`
	ShortURL     string     `json:"short_url,omitempty"`
	Tags         []Tag      `json:"tags,omitempty"`
}

type CreateLinkRequest struct {
	URL       string `json:"url"`
	Title     string `json:"title,omitempty"`
	CustomCode string `json:"custom_code,omitempty"`
	ExpiresIn  int    `json:"expires_in,omitempty"` // days
	MaxClicks  *int   `json:"max_clicks,omitempty"`
	Password   string `json:"password,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

type UpdateLinkRequest struct {
	Title     *string `json:"title,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	ExpiresAt *string `json:"expires_at,omitempty"`
	MaxClicks *int    `json:"max_clicks,omitempty"`
}

type LinkListResponse struct {
	Links      []Link `json:"links"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
