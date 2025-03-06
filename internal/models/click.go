package models

import "time"

type Click struct {
	ID        int       `json:"id"`
	LinkID    int       `json:"link_id"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Referer   string    `json:"referer,omitempty"`
	Country   string    `json:"country,omitempty"`
	City      string    `json:"city,omitempty"`
	Device    string    `json:"device,omitempty"`
	Browser   string    `json:"browser,omitempty"`
	OS        string    `json:"os,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ClickStats struct {
	TotalClicks    int               `json:"total_clicks"`
	UniqueClicks   int               `json:"unique_clicks"`
	ClicksByDay    []DayCount        `json:"clicks_by_day"`
	TopReferrers   []NameCount       `json:"top_referrers"`
	TopCountries   []NameCount       `json:"top_countries"`
	TopBrowsers    []NameCount       `json:"top_browsers"`
	TopDevices     []NameCount       `json:"top_devices"`
	TopOS          []NameCount       `json:"top_os"`
}

type DayCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type NameCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
